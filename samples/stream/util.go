package stream

import (
	"encoding/binary"

	"github.com/reugn/go-streams"
)

// doStream streams data from the outlet to inlet. Blocks by default.
func doStream(outlet streams.Outlet, inlet streams.Inlet) {
	for element := range outlet.Out() {
		inlet.In() <- element
	}

	close(inlet.In())
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
