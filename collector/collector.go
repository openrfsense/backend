package collector

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/openrfsense/backend/database"
	"github.com/openrfsense/backend/database/models"
	"github.com/openrfsense/common/logging"

	"github.com/hamba/avro/v2"
	"github.com/knadh/koanf"
)

//go:embed sample.avsc
var schemasFs embed.FS

var (
	errors   chan error
	quitChan chan bool

	listener net.Listener
	schema   avro.Schema
)

var log = logging.New().
	WithPrefix("collector").
	WithLevel(logging.DebugLevel).
	WithFlags(logging.FlagsDevelopment)

// Initializes an internal TCP listener on the configured port and starts an accept loop
// which waits for TCP packets (in Avro binary format).
func Start(config *koanf.Koanf) error {
	var err error

	listener, err = net.Listen("tcp", fmt.Sprintf(":%d", config.MustInt("collector.port")))
	if err != nil {
		return err
	}

	quitChan = make(chan bool, 1)
	errors = make(chan error, 2)

	// Initialize schema
	schemaBytes, err := schemasFs.ReadFile("sample.avsc")
	if err != nil {
		return err
	}
	schema = avro.MustParse(string(schemaBytes))

	// Start error logger
	go errorLogger(errors)

	// Start the accept loop
	go accept(errors)

	return nil
}

// Gracefully terminates incoming connections and stop the collector.
func Stop() {
	quitChan <- true
}

// Returns a reference to the (buffered) error channel.
func Errors() <-chan error {
	return errors
}

// The actual accept loop.
func accept(errChan chan<- error) {
	wg := sync.WaitGroup{}

	for {
		select {
		case <-quitChan:
			listener.Close()
			wg.Wait()
			return
		default:
		}
		// Block and accept a connection
		conn, err := listener.Accept()
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			errChan <- err
			continue
		}
		// Set an arbitrary 300ms deadline on the current packet
		// err = conn.SetDeadline(time.Now().Add(300 * time.Millisecond))
		// if err != nil {
		// 	errChan <- err
		// }

		// Start worker process
		wg.Add(1)
		go handleRequest(conn, &wg, errors)
	}
}

// Simple TCP request handler which deserializes raw Avro packates into
// database.Sample objects and saves them to the database.
func handleRequest(conn net.Conn, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	packet := bufio.NewReader(conn)
	avroBytes, err := io.ReadAll(packet)
	if err != nil {
		errChan <- err
	}

	s := models.Sample{}
	err = avro.Unmarshal(schema, avroBytes, &s)
	if err != nil {
		errChan <- err
	}

	sql, args, _ := database.Instance().
		Insert("samples").
		Columns(
			"sensor_id",
			"campaign_id",
			"sample_type",
			"time_seconds",
			"time_microseconds",
			"config_antenna_gain",
			"config_antenna_id",
			"config_center_freq",
			"config_est_noise_floor",
			"config_frequency_correction_factor",
			"config_frontend_gain",
			"config_hopping_strategy",
			"config_iq_balance_calibration",
			"config_rf_sync",
			"config_sampling_rate",
			"config_sig_strength_calibration",
			"config_system_sync",
			"config_extra_conf",
			"data",
		).Values(
		s.SensorId,
		s.CampaignId,
		s.SampleType,
		s.SampleTime.Seconds,
		s.SampleTime.Microseconds,
		s.SampleConfig.AntennaGain,
		s.SampleConfig.AntennaId,
		s.SampleConfig.CenterFreq,
		s.SampleConfig.EstNoiseFloor,
		s.SampleConfig.FrequencyCorrectionFactor,
		s.SampleConfig.FrontendGain,
		s.SampleConfig.HoppingStrategy,
		s.SampleConfig.IqBalanceCalibration,
		s.SampleConfig.RfSync,
		s.SampleConfig.SamplingRate,
		s.SampleConfig.SigStrengthCalibration,
		s.SampleConfig.SystemSync,
		s.SampleConfig.ExtraConf,
		s.Data,
	).ToSql()
	log.Debug(sql)
	log.Debug(args)

	err = database.Do(
		context.Background(),
		sql,
		args...,
	)
	if err != nil {
		errChan <- err
	}
}

// Simple consumer which logs error received on a channel.
func errorLogger(errChan <-chan error) {
	for {
		err := <-errChan
		log.Error(err)
	}
}
