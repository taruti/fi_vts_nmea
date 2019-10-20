package main

import (
	"bufio"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"time"
)

type pos struct {
	Time     int64
	Lat, Lon float64
}

type vesselData struct {
	Name string
	Locs []pos
}

type vmap = map[uint32]*vesselData

func getVessel(m vmap, mmsi int64) *vesselData {
	v := m[uint32(mmsi)]
	if v == nil {
		v = &vesselData{}
		m[uint32(mmsi)] = v
	}
	return v
}

func outGpxTo(w io.Writer, m vmap) error {
	w.Write([]byte("<?xml version=\"1.0\"?><gpx>\n"))
	for mmsi, v := range m {
		name := v.Name
		if name == "" {
			name = strconv.Itoa(int(mmsi))
		}

		//		s := "<extensions><opencpn:viz>1</opencpn:viz></extensions>"
		w.Write([]byte("<trk><name>"))
		xml.EscapeText(w, []byte(name))
		w.Write([]byte("</name><trkseg>\n"))
		plat, plon := math.NaN(), math.NaN()
		for _, l := range v.Locs {
			t := time.Unix(0, l.Time*1000000)
			if math.Abs(plat-l.Lat)+2*math.Abs(plon-l.Lon) > 0.5 {
				fmt.Fprint(w, "</trkseg><trkseg>\n")
			}
			plat, plon = l.Lat, l.Lon
			ts := t.Format(time.RFC3339)
			fmt.Fprintf(w, "<trkpt lat=\"%f\" lon=\"%f\"><time>%s</time></trkpt>\n", l.Lat, l.Lon, ts)
		}
		fmt.Fprintf(w, "</trkseg></trk>\n")
	}
	w.Write([]byte("</gpx>"))
	return nil
}

// CsvToGpx reads CSV data produced by this program from stdin
// and produces output to stdout. This needs to read the whole
// input before it produces any output.
func CsvToGpx() error {
	rd := csv.NewReader(bufio.NewReader(os.Stdin))
	rd.ReuseRecord = true
	rd.FieldsPerRecord = -1
	m := map[uint32]*vesselData{}
	for {
		r, err := rd.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		t, err := strconv.ParseInt(r[0], 10, 64)
		if err != nil {
			return err
		}
		mmsi, err := strconv.ParseInt(r[1], 10, 32)
		if err != nil {
			return err
		}
		if r[2] != "" {
			lat, err := strconv.ParseFloat(r[2], 64)
			if err != nil {
				return err
			}
			lon, err := strconv.ParseFloat(r[3], 64)
			if err != nil {
				return err
			}
			v := getVessel(m, mmsi)
			v.Locs = append(v.Locs, pos{Time: t, Lat: lat, Lon: lon})
		} else if r[4] != "" {
			getVessel(m, mmsi).Name = r[4]
		}
	}
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	return outGpxTo(w, m)
}
