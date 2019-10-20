package main

import (
	"fmt"
	"strings"
)

type csvFormatter struct{}

func (csvFormatter) FormatVesselLocation(l *vesselLocation) ([]byte, error) {
	return []byte(fmt.Sprintf("%d,%d,%f,%f\n", l.TimestampExternal, l.MMSI, l.Lat(), l.Lon())), nil
}

func (csvFormatter) FormatVesselMetadata(l *vesselMetadata) ([]byte, error) {
	n := strings.Replace(l.Name, `"`, `""`, -1)
	return []byte(fmt.Sprintf("%d,%d,,,\"%s\"\n", l.Timestamp, l.MMSI, n)), nil
}
