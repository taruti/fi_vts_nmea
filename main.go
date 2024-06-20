package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var proto = flag.String("proto", "udp", "Protocol to use for sending NMEA e.g. udp, tcp, udp4, tcp6.")
var addr = flag.String("address", "127.0.0.1:10110", "Address to send NMEA data.")
var useStdout = flag.Bool("stdout", false, "Send NMEA data to stdout instead of network.")
var verbose = flag.Bool("v", false, "Be more verbose")
var outFormat = flag.String("format", "nmea", "Output format {nmea,csv}")
var serverURL = flag.String("url", "wss://meri.digitraffic.fi:443/mqtt", "Server to use")
var csvToGpx = flag.Bool("csvtogpx", false, "Convert CSV output produced by this program from stdin to a gpx file to stdout.")
var timestamp = flag.Bool("time", false, "Emit timestamps into output stream")

func openOutput() (io.WriteCloser, error) {
	if *useStdout {
		log.Print("Using NMEA output stdout")
		return os.Stdout, nil
	}
	log.Printf("Using NMEA output %q %q", *proto, *addr)
	return net.Dial(*proto, *addr)
}

func debugf(s string, vs ...interface{}) {
	if *verbose {
		log.Printf(s, vs...)
	}
}
func debug(vs ...interface{}) {
	if *verbose {
		log.Print(vs...)
	}
}

type outFormatter interface {
	FormatVesselLocation(*vesselLocation) ([]byte, error)
	FormatVesselMetadata(*vesselMetadata) ([]byte, error)
	FormatTime(time.Time) []byte
}

func getFormatter() outFormatter {
	switch *outFormat {
	case "csv":
		return csvFormatter{}
	case "nmea":
		return nmeaFormatter{}
	}
	log.Fatalf("Invalid output format: %q", *outFormat)
	return nil // not reachable
}

func work() error {
	if *csvToGpx {
		return CsvToGpx()
	}

	out, err := openOutput()
	if err != nil {
		return err
	}
	defer out.Close()

	w := &cc{out, getFormatter(), make(chan error, 1)}
	err = dialMqtt(w.onMessageReceived)
	if err != nil {
		return err
	}

	return <-w.errCh
}

func main() {
	flag.Parse()

	err := work()
	if err != nil {
		log.Println(err)
	}
}
