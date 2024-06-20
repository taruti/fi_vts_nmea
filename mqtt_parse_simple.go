package main

import (
	"encoding/json"
)

// MMSI is a unique identifier for e.g. ships, see
// https://en.wikipedia.org/wiki/Maritime_Mobile_Service_Identity
type MMSI uint32

type vesselLocation struct {
	MMSI      MMSI `json:"MMSI"`
	Timestamp int  `json:"time"`
	///	Type             string   `json:"type"`
	Geometry geometry `json:"geometry"`
	//	vesselProperties `json:"properties"`
	Sog     float64 `json:"sog"`
	Cog     float64 `json:"cog"`
	NavStat int     `json:"navStat"`
	Rot     int     `json:"rot"`
	PosAcc  bool    `json:"posAcc"`
	Raim    bool    `json:"raim"`
	Heading int     `json:"heading"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

/*
"vessels-v2/230010760/location"
Message: {
  "time" : 1686936661,
  "sog" : 7.4,
  "cog" : 290.6,
  "navStat" : 0,
  "rot" : 0,
  "posAcc" : false,
  "raim" : false,
  "heading" : 286,
  "lon" : 28.460167,
  "lat" : 62.101388
}
*/

type vesselProperties struct {
	Timestamp         int   `json:"timestamp"`
	TimestampExternal int64 `json:"timestampExternal"`
}
type geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type vesselMetadata struct {
	Timestamp       int64  `json:"timestamp"`
	Destination     string `json:"destination"`
	MMSI            MMSI   `json:"mmsi"`
	ShipType        int    `json:"shipType"`
	CallSign        string `json:"callSign"`
	IMO             int    `json:"imo"`
	Draught         int    `json:"draught"`
	ETA             int    `json:"eta"`
	PosType         int    `json:"posType"`
	ReferencePointA int    `json:"referencePointA"`
	ReferencePointB int    `json:"referencePointB"`
	ReferencePointC int    `json:"referencePointC"`
	ReferencePointD int    `json:"referencePointD"`
	Name            string `json:"name"`
}

func parseVesselLocation(bs []byte, msg *vesselLocation) error {
	return json.Unmarshal(bs, msg)
}
func parseVesselMetadata(bs []byte, msg *vesselMetadata) error {
	return json.Unmarshal(bs, msg)
}
