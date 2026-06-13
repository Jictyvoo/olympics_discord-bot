package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/fnv"
)

func Hash[T any](object T) (string, error) {
	h := fnv.New64a() // Create a new FNV-1a 64-bit hash instance

	// Serialize POSMessageConfig into gob
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(object); err != nil {
		return "", err
	}

	// Write the serialized bytes to the hash
	if _, err := h.Write(buf.Bytes()); err != nil {
		return "", err
	}

	// Return the resulting hash value as a hexadecimal string
	return fmt.Sprintf("%x", h.Sum64()), nil
}
