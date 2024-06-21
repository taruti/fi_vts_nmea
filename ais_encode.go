package main

import (
	"github.com/taruti/bstream"

	"fmt"
	"time"
)

// see http://catb.org/gpsd/AIVDM.html

func i4m(f float64) uint64 {
	return uint64(int64((600000 * f) + 0.5))
}
func u1(f float64) uint64 {
	return uint64((10 * f) + 0.5)
}

type nmeaFormatter struct{}

func (nmeaFormatter) FormatVesselLocation(l *vesselLocation) ([]byte, error) {
	const lenBits = 168
	b := bstream.NewBStreamWriter(lenBits / 8)
	b.WriteBits(1, 6)
	b.WriteBits(3, 2)
	b.WriteBits(uint64(l.MMSI), 30)
	b.WriteBits(uint64(l.NavStat), 4)
	b.WriteBits(uint64(l.Rot), 8)
	b.WriteBits(u1(l.Sog), 10)
	b.WriteBool(l.PosAcc)
	b.WriteBits(i4m(l.Lon), 28)
	b.WriteBits(i4m(l.Lat), 27)
	b.WriteBits(u1(l.Cog), 12)
	b.WriteBits(uint64(l.Heading), 9)
	b.WriteBits(uint64(l.Timestamp), 6)
	b.WriteBits(0, 2+3) // maneuver indicator  + spare
	b.WriteBool(l.Raim)
	b.WriteBits(0, 19)
	bs := make([]byte, 0, 82)
	bs = append(bs, `!AIVDM,1,1,,A,`...)
	bs = nmeaAISAppend(bs, b, lenBits/6)
	bs = append(bs, `,0`...)
	bs = nmeaAppendChecksum(bs)
	return bs, nil
}

func encodeChar(b byte) byte {
	switch {
	case b >= '@' && b <= '_':
		return b - '@'
	case b >= ' ' && b <= '?':
		return b
	default:
		return '?'
	}
}

func writeString(dst *bstream.BStream, toEncode string, length int) {
	for i := 0; i < length; i++ {
		b := byte('@')
		if i < len(toEncode) {
			b = encodeChar(toEncode[i])
		}
		dst.WriteBits(uint64(b), 6)
	}
}

func (nmeaFormatter) FormatVesselMetadata(l *vesselMetadata) ([]byte, error) {
	const lenBits = 424
	b := bstream.NewBStreamWriter(lenBits / 8)
	b.WriteBits(5, 6)
	b.WriteBits(3, 2)
	b.WriteBits(uint64(l.MMSI), 30)
	b.WriteBits(0, 2)
	b.WriteBits(uint64(l.IMO), 30)
	writeString(b, l.CallSign, 7)
	writeString(b, l.Name, 20)
	b.WriteBits(uint64(l.ShipType), 8)
	b.WriteBits(uint64(l.ReferencePointA), 9)
	b.WriteBits(uint64(l.ReferencePointB), 9)
	b.WriteBits(uint64(l.ReferencePointC), 6)
	b.WriteBits(uint64(l.ReferencePointD), 6)
	b.WriteBits(uint64(l.PosType), 4)
	b.WriteBits(uint64(l.ETA), 20)
	b.WriteBits(uint64(l.Draught), 8)
	writeString(b, l.Destination, 20)
	b.WriteBits(0, 2+2) // 2 extra bits to make it 6-compatible.
	bs := make([]byte, 0, 82*2)
	bs = append(bs, `!AIVDM,2,1,,A,`...)
	bs = nmeaAISAppend(bs, b, 35)
	bs = append(bs, `,0`...)
	bs = nmeaAppendChecksum(bs)
	offset := len(bs) + 1
	bs = append(bs, `!AIVDM,2,2,,A,`...)
	bs = nmeaAISAppend(bs, b, 36)
	bs = append(bs, `,2`...)
	bs = nmeaAppendChecksumFrom(offset, bs)
	debugf("NMEA: %s", bs)
	return bs, nil
}

func (nmeaFormatter) FormatTime(t time.Time) []byte {
	// FIXME
	bs := []byte(fmt.Sprintf("\\c:%d*00\\$ZCRTE,1,1,c,0*00\n", t.Unix()))
	debugf("NMEA: %s", bs)
	return bs
}

var encoderChars = [64]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	':', ';', '<', '=', '>', '?', '@',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W',
	'`',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w',
}

func nmeaAISAppend(bs []byte, src *bstream.BStream, length int) []byte {
	for i := 0; i < length; i++ {
		v, err := src.ReadBits(6)
		if err != nil {
			panic("Encoding error: " + err.Error())
		}
		bs = append(bs, encoderChars[int(v)])
	}
	return bs
}

func nmeaAppendChecksumFrom(offset int, bs []byte) []byte {
	var csum byte
	for _, by := range bs[offset:] {
		csum ^= by
	}
	bs = append(bs, '*')
	bs = append(bs, hexChars[csum>>4])
	bs = append(bs, hexChars[csum&0xF])
	bs = append(bs, "\r\n"...)
	return bs
}

func nmeaAppendChecksum(bs []byte) []byte {
	return nmeaAppendChecksumFrom(1, bs)
}

var hexChars = [16]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'A', 'B', 'C', 'D', 'E', 'F',
}
