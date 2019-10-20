# Finnish AIS information from VTS as NMEA stream

+ **Not for navigation or safety purposes**, the data may be wrong, old, and misintepreted!
+ **This code is not endorsed by the VTS Finland/TMFG/Väylä/Traficom/...**
+ Only class A AIS data, VTS filters out class B vessels.
+ Code is MIT licensed, data stream is owned by VTS Finland.
+ Also supports CSV output of VTS data stream.

## Building from Source

+ Install Go https://golang.org
+ Run `go get github.com/taruti/fi_vts_nmea`
+ Now you have installed `fi_vts_nmea`

## Help

```
Usage of ./fi_vts_nmea:
  -address string
    	Address to send NMEA data. (default "127.0.0.1:10110")
  -csvtogpx
    	Convert CSV output produced by this program from stdin to a gpx file to stdout.
  -format string
    	Output format {nmea,csv} (default "nmea")
  -proto string
    	Protocol to use for sending NMEA e.g. udp, tcp, udp4, tcp6. (default "udp")
  -server string
    	Server to use (default "meri-test.digitraffic.fi")
  -stdout
    	Send NMEA data to stdout instead of network.
  -v	Be more verbose
```

