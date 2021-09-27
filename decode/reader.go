package decode

import (
	"errors"
)

// ErrorOutOfIndex returned if try to read bytes
// more than reader contains.
var ErrorOutOfIndex = errors.New("array out of index")

type reader struct {
	pos     int
	payload []byte
}

func (r *reader) getBytes(n int) ([]byte, error) {
	if n+r.pos > len(r.payload) {
		return nil, ErrorOutOfIndex
	}
	payload := r.payload[r.pos : r.pos+n]
	r.pos += n
	return payload, nil
}
