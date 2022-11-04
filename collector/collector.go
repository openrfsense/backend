package collector

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/knadh/koanf"
	"github.com/linkedin/goavro/v2"
	"github.com/openrfsense/backend/database"
	"github.com/openrfsense/common/logging"
)

//go:embed sample.avsc
var schemasFs embed.FS

var (
	listener *net.TCPListener
	quitChan chan bool

	codec *goavro.Codec
)

var log = logging.New().
	WithPrefix("collector").
	WithLevel(logging.DebugLevel).
	WithFlags(logging.FlagsDevelopment)

// Initializes an internal TCP listener on the configured port and starts an accept loop
// which waits for TCP packets (in Avro binary format).
func Start(config *koanf.Koanf) error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", config.MustInt("collector.port")))
	if err != nil {
		return err
	}
	listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	quitChan = make(chan bool, 1)

	codecBytes, err := schemasFs.ReadFile("sample.avsc")
	if err != nil {
		return err
	}
	codec, err = goavro.NewCodec(string(codecBytes))
	if err != nil {
		return err
	}

	go accept()

	return nil
}

// Gracefully terminates incoming connections and stop the collector.
func Stop() {
	quitChan <- true
}

// The actual accept loop.
func accept() {
	wg := sync.WaitGroup{}

	for {
		select {
		case <-quitChan:
			listener.Close()
			wg.Wait()
			return
		default:
		}
		err := listener.SetDeadline(time.Now().Add(1e9))
		if err != nil {
			log.Errorf("Deadline: %v", err)
			continue
		}
		conn, err := listener.AcceptTCP()
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			continue
		}
		if err != nil {
			log.Errorf("AcceptTCP: %v", err)
			continue
		}
		wg.Add(1)
		go func() {
			wg.Done()
			err = handleRequest(conn)
			if err != nil {
				log.Error(err)
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

	native, _, err := codec.NativeFromBinary(avroBytes)
	if err != nil {
		return err
	}

	textual, err := codec.TextualFromNative(nil, native)
	if err != nil {
		return err
	}

	sample := database.Sample{}
	err = json.Unmarshal(textual, &sample)
	if err != nil {
		return err
	}
	return database.Instance().Create(&sample).Error
}
