package store

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register(&CachedResponse{})
}

func serialize(value interface{}) ([]byte, error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func unserialize(payload []byte, ptr interface{}) (err error) {
	return gob.NewDecoder(bytes.NewBuffer(payload)).Decode(ptr)
}
