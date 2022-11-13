package collector

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/openrfsense/backend/database"
	"github.com/openrfsense/common/logging"
	"github.com/valyala/tcplisten"

	"github.com/hamba/avro/v2"
	"github.com/knadh/koanf"
)

//go:embed sample.avsc
var schemasFs embed.FS

var (
	listener net.Listener
	quitChan chan bool
	errors   chan error

	schema avro.Schema
)

var log = logging.New().
	WithPrefix("collector").
	WithLevel(logging.DebugLevel).
	WithFlags(logging.FlagsDevelopment)

// Initializes an internal TCP listener on the configured port and starts an accept loop
// which waits for TCP packets (in Avro binary format).
func Start(config *koanf.Koanf) error {
	var err error

	conf := &tcplisten.Config{}
	listener, err = conf.NewListener("tcp4", fmt.Sprintf(":%d", config.MustInt("collector.port")))
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
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			errChan <- err
			continue
		}
		// Set an arbitrary 300ms deadline on the current packet
		err = conn.SetDeadline(time.Now().Add(300 * time.Millisecond))
		if err != nil {
			errChan <- err
		}

		// Start worker process
		wg.Add(1)
		go func() {
			wg.Done()
			err = handleRequest(conn)
			if err != nil {
				errChan <- err
			}
		}()
	}
}

// Simple TCP request handler which deserializes raw Avro packates into
// database.Sample objects and saves them to the database.
func handleRequest(conn net.Conn) error {
	packet := bufio.NewReader(conn)
	avroBytes, err := io.ReadAll(packet)
	if err != nil {
		return err
	}

	sample := database.Sample{}
	err = avro.Unmarshal(schema, avroBytes, &sample)
	if err != nil {
		return err
	}

	return database.Instance().Create(&sample).Error
}

// Simple consumer which logs error received on a channel.
func errorLogger(errChan <-chan error) {
	for {
		err := <-errChan
		log.Error(err)
	}
}
