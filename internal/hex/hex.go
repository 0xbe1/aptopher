// Package hex provides utilities for hex encoding/decoding with 0x prefix.
package hex

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// Encode encodes bytes to a hex string with 0x prefix.
func Encode(data []byte) string {
	return "0x" + hex.EncodeToString(data)
}

// Decode decodes a hex string (with or without 0x prefix) to bytes.
// Handles odd-length hex strings by left-padding with a zero.
func Decode(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	// Pad odd-length strings with a leading zero
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}

// MustDecode decodes a hex string or panics.
func MustDecode(s string) []byte {
	data, err := Decode(s)
	if err != nil {
		panic(fmt.Sprintf("hex: invalid hex string %q: %v", s, err))
	}
	return data
}

// DecodeFixed decodes a hex string to a fixed-size byte array.
// Returns error if the decoded length doesn't match expected size.
func DecodeFixed(s string, size int) ([]byte, error) {
	data, err := Decode(s)
	if err != nil {
		return nil, err
	}
	if len(data) > size {
		return nil, fmt.Errorf("hex: decoded length %d exceeds expected size %d", len(data), size)
	}
	// Left-pad with zeros if shorter
	if len(data) < size {
		padded := make([]byte, size)
		copy(padded[size-len(data):], data)
		return padded, nil
	}
	return data, nil
}
