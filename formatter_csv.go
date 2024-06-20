package main

import (
	"fmt"
	"strings"
	"time"
)

type csvFormatter struct{}

func (csvFormatter) FormatVesselLocation(l *vesselLocation) ([]byte, error) {
	return []byte(fmt.Sprintf("%d,%d,%f,%f\n", l.Timestamp, l.MMSI, l.Lat, l.Lon)), nil
}

func (csvFormatter) FormatVesselMetadata(l *vesselMetadata) ([]byte, error) {
	n := strings.Replace(l.Name, `"`, `""`, -1)
	return []byte(fmt.Sprintf("%d,%d,,,\"%s\"\n", l.Timestamp, l.MMSI, n)), nil
}

func (csvFormatter) FormatTime(t time.Time) []byte {
	bs := []byte(fmt.Sprintf("t,%d\n", t.Unix()))
	return bs
}
