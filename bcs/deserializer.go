package bcs

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"
)

// deserializerPool provides pooled Deserializer instances to reduce allocations.
var deserializerPool = sync.Pool{
	New: func() interface{} { return &Deserializer{} },
}

// AcquireDeserializer returns a Deserializer from the pool, initialized with the given data.
// Call ReleaseDeserializer when done to return it to the pool.
func AcquireDeserializer(data []byte) *Deserializer {
	d := deserializerPool.Get().(*Deserializer)
	d.data = data
	d.offset = 0
	d.err = nil
	return d
}

// ReleaseDeserializer returns a Deserializer to the pool.
// Do not use the Deserializer after releasing it.
func ReleaseDeserializer(d *Deserializer) {
	d.data = nil // Allow GC
	deserializerPool.Put(d)
}

// Deserializer provides BCS deserialization capabilities.
type Deserializer struct {
	data   []byte
	offset int
	err    error
}

// NewDeserializer creates a new BCS deserializer from the given bytes.
func NewDeserializer(data []byte) *Deserializer {
	return &Deserializer{data: data}
}

// Error returns any error that occurred during deserialization.
func (d *Deserializer) Error() error {
	return d.err
}

// SetError sets an error on the deserializer. Once set, subsequent operations are no-ops.
func (d *Deserializer) SetError(err error) {
	if d.err == nil {
		d.err = err
	}
}

// Remaining returns the number of bytes remaining to be read.
func (d *Deserializer) Remaining() int {
	return len(d.data) - d.offset
}

// checkRemaining verifies there are at least n bytes remaining.
func (d *Deserializer) checkRemaining(n int) bool {
	if d.err != nil {
		return false
	}
	if d.offset+n > len(d.data) {
		d.SetError(fmt.Errorf("bcs: unexpected end of input, need %d bytes, have %d", n, d.Remaining()))
		return false
	}
	return true
}

// Bool deserializes a boolean value.
func (d *Deserializer) Bool() bool {
	if !d.checkRemaining(1) {
		return false
	}
	b := d.data[d.offset]
	d.offset++
	switch b {
	case 0x00:
		return false
	case 0x01:
		return true
	default:
		d.SetError(fmt.Errorf("bcs: invalid boolean value: 0x%02x", b))
		return false
	}
}

// U8 deserializes an unsigned 8-bit integer.
func (d *Deserializer) U8() uint8 {
	if !d.checkRemaining(1) {
		return 0
	}
	v := d.data[d.offset]
	d.offset++
	return v
}

// U16 deserializes an unsigned 16-bit integer in little-endian format.
func (d *Deserializer) U16() uint16 {
	if !d.checkRemaining(2) {
		return 0
	}
	v := binary.LittleEndian.Uint16(d.data[d.offset:])
	d.offset += 2
	return v
}

// U32 deserializes an unsigned 32-bit integer in little-endian format.
func (d *Deserializer) U32() uint32 {
	if !d.checkRemaining(4) {
		return 0
	}
	v := binary.LittleEndian.Uint32(d.data[d.offset:])
	d.offset += 4
	return v
}

// U64 deserializes an unsigned 64-bit integer in little-endian format.
func (d *Deserializer) U64() uint64 {
	if !d.checkRemaining(8) {
		return 0
	}
	v := binary.LittleEndian.Uint64(d.data[d.offset:])
	d.offset += 8
	return v
}

// U128 deserializes a 128-bit unsigned integer in little-endian format.
func (d *Deserializer) U128() *big.Int {
	if !d.checkRemaining(16) {
		return nil
	}
	// Use stack array instead of allocating slice
	var buf [16]byte
	// Copy and reverse in one pass (little-endian to big-endian)
	for i := 0; i < 16; i++ {
		buf[15-i] = d.data[d.offset+i]
	}
	d.offset += 16
	return new(big.Int).SetBytes(buf[:])
}

// U256 deserializes a 256-bit unsigned integer in little-endian format.
func (d *Deserializer) U256() *big.Int {
	if !d.checkRemaining(32) {
		return nil
	}
	// Use stack array instead of allocating slice
	var buf [32]byte
	// Copy and reverse in one pass (little-endian to big-endian)
	for i := 0; i < 32; i++ {
		buf[31-i] = d.data[d.offset+i]
	}
	d.offset += 32
	return new(big.Int).SetBytes(buf[:])
}

// Uleb128 deserializes an unsigned integer using ULEB128 variable-length encoding.
func (d *Deserializer) Uleb128() uint32 {
	if d.err != nil {
		return 0
	}
	var result uint32
	var shift uint
	for {
		if !d.checkRemaining(1) {
			return 0
		}
		b := d.data[d.offset]
		d.offset++
		result |= uint32(b&0x7f) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift > 28 {
			d.SetError(fmt.Errorf("bcs: ULEB128 overflow"))
			return 0
		}
	}
	return result
}

// Bytes deserializes a byte slice with a ULEB128 length prefix.
func (d *Deserializer) Bytes() []byte {
	length := d.Uleb128()
	if d.err != nil {
		return nil
	}
	if !d.checkRemaining(int(length)) {
		return nil
	}
	result := make([]byte, length)
	copy(result, d.data[d.offset:d.offset+int(length)])
	d.offset += int(length)
	return result
}

// FixedBytes deserializes a fixed-size byte slice without a length prefix.
func (d *Deserializer) FixedBytes(length int) []byte {
	if !d.checkRemaining(length) {
		return nil
	}
	result := make([]byte, length)
	copy(result, d.data[d.offset:d.offset+length])
	d.offset += length
	return result
}

// FixedBytesNoCopy returns a slice of the internal buffer without copying.
// WARNING: Do not modify the returned slice. Do not use it after the
// Deserializer is released back to the pool or the underlying data is modified.
// Use this only in performance-critical paths where the data is immediately consumed.
func (d *Deserializer) FixedBytesNoCopy(length int) []byte {
	if !d.checkRemaining(length) {
		return nil
	}
	result := d.data[d.offset : d.offset+length]
	d.offset += length
	return result
}

// String deserializes a UTF-8 string with a ULEB128 length prefix.
func (d *Deserializer) String() string {
	return string(d.Bytes())
}

// Struct deserializes a type that implements Unmarshaler.
func (d *Deserializer) Struct(v Unmarshaler) {
	if d.err != nil {
		return
	}
	v.UnmarshalBCS(d)
}

// DeserializeSequence deserializes a slice of Unmarshaler elements.
// Reads ULEB128 length followed by each element.
func DeserializeSequence[T Unmarshaler](d *Deserializer, factory func() T) []T {
	if d.err != nil {
		return nil
	}
	length := d.Uleb128()
	if d.err != nil {
		return nil
	}
	result := make([]T, length)
	for i := uint32(0); i < length; i++ {
		item := factory()
		item.UnmarshalBCS(d)
		if d.err != nil {
			return nil
		}
		result[i] = item
	}
	return result
}

// DeserializeOption deserializes an optional value.
// Reads 0x00 for nil (None) or 0x01 followed by the value (Some).
func DeserializeOption[T Unmarshaler](d *Deserializer, factory func() T) *T {
	if d.err != nil {
		return nil
	}
	hasValue := d.U8()
	if d.err != nil {
		return nil
	}
	switch hasValue {
	case 0:
		return nil
	case 1:
		item := factory()
		item.UnmarshalBCS(d)
		if d.err != nil {
			return nil
		}
		return &item
	default:
		d.SetError(fmt.Errorf("bcs: invalid option tag: 0x%02x", hasValue))
		return nil
	}
}

// Deserialize is a convenience function to deserialize bytes into an Unmarshaler.
func Deserialize[T Unmarshaler](data []byte, v T) error {
	d := NewDeserializer(data)
	v.UnmarshalBCS(d)
	if d.err != nil {
		return d.err
	}
	if d.Remaining() > 0 {
		return fmt.Errorf("bcs: %d bytes remaining after deserialization", d.Remaining())
	}
	return nil
}
