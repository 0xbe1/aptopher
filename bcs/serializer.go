package bcs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
)

// Serializer provides BCS serialization capabilities.
// It uses a streaming interface to minimize allocations.
type Serializer struct {
	buf bytes.Buffer
	err error
}

// NewSerializer creates a new BCS serializer.
func NewSerializer() *Serializer {
	return &Serializer{}
}

// Error returns any error that occurred during serialization.
func (s *Serializer) Error() error {
	return s.err
}

// SetError sets an error on the serializer. Once set, subsequent operations are no-ops.
func (s *Serializer) SetError(err error) {
	if s.err == nil {
		s.err = err
	}
}

// ToBytes returns the serialized bytes. Returns nil if an error occurred.
func (s *Serializer) ToBytes() []byte {
	if s.err != nil {
		return nil
	}
	return s.buf.Bytes()
}

// Bool serializes a boolean value.
// BCS: 0x00 for false, 0x01 for true
func (s *Serializer) Bool(v bool) {
	if s.err != nil {
		return
	}
	if v {
		s.buf.WriteByte(0x01)
	} else {
		s.buf.WriteByte(0x00)
	}
}

// U8 serializes an unsigned 8-bit integer.
func (s *Serializer) U8(v uint8) {
	if s.err != nil {
		return
	}
	s.buf.WriteByte(v)
}

// U16 serializes an unsigned 16-bit integer in little-endian format.
func (s *Serializer) U16(v uint16) {
	if s.err != nil {
		return
	}
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], v)
	s.buf.Write(buf[:])
}

// U32 serializes an unsigned 32-bit integer in little-endian format.
func (s *Serializer) U32(v uint32) {
	if s.err != nil {
		return
	}
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], v)
	s.buf.Write(buf[:])
}

// U64 serializes an unsigned 64-bit integer in little-endian format.
func (s *Serializer) U64(v uint64) {
	if s.err != nil {
		return
	}
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], v)
	s.buf.Write(buf[:])
}

// U128 serializes a 128-bit unsigned integer in little-endian format.
func (s *Serializer) U128(v *big.Int) {
	if s.err != nil {
		return
	}
	if v == nil {
		s.SetError(fmt.Errorf("bcs: U128 value is nil"))
		return
	}
	if v.Sign() < 0 {
		s.SetError(fmt.Errorf("bcs: U128 value is negative"))
		return
	}
	// U128 is 16 bytes in little-endian
	bytes := v.Bytes() // big-endian
	if len(bytes) > 16 {
		s.SetError(fmt.Errorf("bcs: U128 value too large"))
		return
	}
	// Pad to 16 bytes and reverse for little-endian
	var buf [16]byte
	for i, b := range bytes {
		buf[len(bytes)-1-i] = b
	}
	s.buf.Write(buf[:])
}

// U256 serializes a 256-bit unsigned integer in little-endian format.
func (s *Serializer) U256(v *big.Int) {
	if s.err != nil {
		return
	}
	if v == nil {
		s.SetError(fmt.Errorf("bcs: U256 value is nil"))
		return
	}
	if v.Sign() < 0 {
		s.SetError(fmt.Errorf("bcs: U256 value is negative"))
		return
	}
	// U256 is 32 bytes in little-endian
	bytes := v.Bytes() // big-endian
	if len(bytes) > 32 {
		s.SetError(fmt.Errorf("bcs: U256 value too large"))
		return
	}
	// Pad to 32 bytes and reverse for little-endian
	var buf [32]byte
	for i, b := range bytes {
		buf[len(bytes)-1-i] = b
	}
	s.buf.Write(buf[:])
}

// Uleb128 serializes an unsigned integer using ULEB128 variable-length encoding.
// Used for sequence lengths and enum variants.
func (s *Serializer) Uleb128(v uint32) {
	if s.err != nil {
		return
	}
	for {
		b := byte(v & 0x7f)
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		s.buf.WriteByte(b)
		if v == 0 {
			break
		}
	}
}

// Bytes serializes a byte slice with a ULEB128 length prefix.
func (s *Serializer) Bytes(v []byte) {
	if s.err != nil {
		return
	}
	s.Uleb128(uint32(len(v)))
	s.buf.Write(v)
}

// FixedBytes serializes a byte slice without a length prefix.
// Use for fixed-size types like AccountAddress.
func (s *Serializer) FixedBytes(v []byte) {
	if s.err != nil {
		return
	}
	s.buf.Write(v)
}

// String serializes a UTF-8 string with a ULEB128 length prefix.
func (s *Serializer) String(v string) {
	s.Bytes([]byte(v))
}

// Struct serializes a type that implements Marshaler.
func (s *Serializer) Struct(v Marshaler) {
	if s.err != nil {
		return
	}
	v.MarshalBCS(s)
}

// SerializeSequence serializes a slice of Marshaler elements.
// Writes ULEB128 length followed by each element.
func SerializeSequence[T Marshaler](s *Serializer, items []T) {
	if s.err != nil {
		return
	}
	s.Uleb128(uint32(len(items)))
	for _, item := range items {
		item.MarshalBCS(s)
	}
}

// SerializeOption serializes an optional value.
// Writes 0x00 for nil (None) or 0x01 followed by the value (Some).
func SerializeOption[T Marshaler](s *Serializer, opt *T) {
	if s.err != nil {
		return
	}
	if opt == nil {
		s.U8(0)
	} else {
		s.U8(1)
		(*opt).MarshalBCS(s)
	}
}

// Serialize is a convenience function to serialize a Marshaler to bytes.
func Serialize(v Marshaler) ([]byte, error) {
	s := NewSerializer()
	v.MarshalBCS(s)
	if s.err != nil {
		return nil, s.err
	}
	return s.ToBytes(), nil
}

// SerializeU8 serializes a uint8 to bytes.
func SerializeU8(v uint8) []byte {
	return []byte{v}
}

// SerializeU64 serializes a uint64 to bytes.
func SerializeU64(v uint64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], v)
	return buf[:]
}

// SerializeString serializes a string to bytes.
func SerializeString(v string) []byte {
	s := NewSerializer()
	s.String(v)
	return s.ToBytes()
}
