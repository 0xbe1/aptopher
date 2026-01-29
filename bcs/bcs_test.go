package bcs

import (
	"bytes"
	"math/big"
	"testing"
)

func TestBool(t *testing.T) {
	tests := []struct {
		name  string
		value bool
		want  []byte
	}{
		{"false", false, []byte{0x00}},
		{"true", true, []byte{0x01}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSerializer()
			s.Bool(tt.value)
			if !bytes.Equal(s.ToBytes(), tt.want) {
				t.Errorf("Bool(%v) = %v, want %v", tt.value, s.ToBytes(), tt.want)
			}
			// Roundtrip
			d := NewDeserializer(tt.want)
			got := d.Bool()
			if d.Error() != nil {
				t.Errorf("Bool deserialize error: %v", d.Error())
			}
			if got != tt.value {
				t.Errorf("Bool roundtrip = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestU8(t *testing.T) {
	tests := []struct {
		value uint8
		want  []byte
	}{
		{0, []byte{0x00}},
		{127, []byte{0x7f}},
		{255, []byte{0xff}},
	}
	for _, tt := range tests {
		s := NewSerializer()
		s.U8(tt.value)
		if !bytes.Equal(s.ToBytes(), tt.want) {
			t.Errorf("U8(%v) = %v, want %v", tt.value, s.ToBytes(), tt.want)
		}
		d := NewDeserializer(tt.want)
		got := d.U8()
		if got != tt.value {
			t.Errorf("U8 roundtrip = %v, want %v", got, tt.value)
		}
	}
}

func TestU16(t *testing.T) {
	tests := []struct {
		value uint16
		want  []byte
	}{
		{0, []byte{0x00, 0x00}},
		{256, []byte{0x00, 0x01}},
		{0x1234, []byte{0x34, 0x12}}, // little-endian
	}
	for _, tt := range tests {
		s := NewSerializer()
		s.U16(tt.value)
		if !bytes.Equal(s.ToBytes(), tt.want) {
			t.Errorf("U16(%v) = %v, want %v", tt.value, s.ToBytes(), tt.want)
		}
		d := NewDeserializer(tt.want)
		got := d.U16()
		if got != tt.value {
			t.Errorf("U16 roundtrip = %v, want %v", got, tt.value)
		}
	}
}

func TestU32(t *testing.T) {
	tests := []struct {
		value uint32
		want  []byte
	}{
		{0, []byte{0x00, 0x00, 0x00, 0x00}},
		{0x12345678, []byte{0x78, 0x56, 0x34, 0x12}}, // little-endian
	}
	for _, tt := range tests {
		s := NewSerializer()
		s.U32(tt.value)
		if !bytes.Equal(s.ToBytes(), tt.want) {
			t.Errorf("U32(%v) = %v, want %v", tt.value, s.ToBytes(), tt.want)
		}
		d := NewDeserializer(tt.want)
		got := d.U32()
		if got != tt.value {
			t.Errorf("U32 roundtrip = %v, want %v", got, tt.value)
		}
	}
}

func TestU64(t *testing.T) {
	tests := []struct {
		value uint64
		want  []byte
	}{
		{0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{0x123456789abcdef0, []byte{0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12}},
	}
	for _, tt := range tests {
		s := NewSerializer()
		s.U64(tt.value)
		if !bytes.Equal(s.ToBytes(), tt.want) {
			t.Errorf("U64(%v) = %v, want %v", tt.value, s.ToBytes(), tt.want)
		}
		d := NewDeserializer(tt.want)
		got := d.U64()
		if got != tt.value {
			t.Errorf("U64 roundtrip = %v, want %v", got, tt.value)
		}
	}
}

func TestUleb128(t *testing.T) {
	tests := []struct {
		value uint32
		want  []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{127, []byte{0x7f}},
		{128, []byte{0x80, 0x01}},
		{16384, []byte{0x80, 0x80, 0x01}},
	}
	for _, tt := range tests {
		s := NewSerializer()
		s.Uleb128(tt.value)
		if !bytes.Equal(s.ToBytes(), tt.want) {
			t.Errorf("Uleb128(%v) = %v, want %v", tt.value, s.ToBytes(), tt.want)
		}
		d := NewDeserializer(tt.want)
		got := d.Uleb128()
		if got != tt.value {
			t.Errorf("Uleb128 roundtrip = %v, want %v", got, tt.value)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		value string
		want  []byte
	}{
		{"", []byte{0x00}},
		{"hello", []byte{0x05, 'h', 'e', 'l', 'l', 'o'}},
		{"世界", []byte{0x06, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c}}, // UTF-8
	}
	for _, tt := range tests {
		s := NewSerializer()
		s.String(tt.value)
		if !bytes.Equal(s.ToBytes(), tt.want) {
			t.Errorf("String(%q) = %v, want %v", tt.value, s.ToBytes(), tt.want)
		}
		d := NewDeserializer(tt.want)
		got := d.String()
		if got != tt.value {
			t.Errorf("String roundtrip = %q, want %q", got, tt.value)
		}
	}
}

func TestBytes(t *testing.T) {
	tests := []struct {
		value []byte
		want  []byte
	}{
		{[]byte{}, []byte{0x00}},
		{[]byte{0x01, 0x02, 0x03}, []byte{0x03, 0x01, 0x02, 0x03}},
	}
	for _, tt := range tests {
		s := NewSerializer()
		s.Bytes(tt.value)
		if !bytes.Equal(s.ToBytes(), tt.want) {
			t.Errorf("Bytes(%v) = %v, want %v", tt.value, s.ToBytes(), tt.want)
		}
		d := NewDeserializer(tt.want)
		got := d.Bytes()
		if !bytes.Equal(got, tt.value) {
			t.Errorf("Bytes roundtrip = %v, want %v", got, tt.value)
		}
	}
}

func TestFixedBytes(t *testing.T) {
	value := []byte{0x01, 0x02, 0x03, 0x04}
	s := NewSerializer()
	s.FixedBytes(value)
	if !bytes.Equal(s.ToBytes(), value) {
		t.Errorf("FixedBytes(%v) = %v, want %v", value, s.ToBytes(), value)
	}
	d := NewDeserializer(value)
	got := d.FixedBytes(4)
	if !bytes.Equal(got, value) {
		t.Errorf("FixedBytes roundtrip = %v, want %v", got, value)
	}
}

func TestFixedBytesNoCopy(t *testing.T) {
	value := []byte{0x01, 0x02, 0x03, 0x04}
	d := NewDeserializer(value)
	got := d.FixedBytesNoCopy(4)
	if !bytes.Equal(got, value) {
		t.Errorf("FixedBytesNoCopy = %v, want %v", got, value)
	}
	if d.Remaining() != 0 {
		t.Errorf("FixedBytesNoCopy should consume all bytes, remaining: %d", d.Remaining())
	}
	// Verify it's actually a slice of the original (shares underlying array)
	got[0] = 0xff
	if value[0] != 0xff {
		t.Error("FixedBytesNoCopy should return slice of original data")
	}
}

func TestBytesNoCopy(t *testing.T) {
	// Serialize bytes with length prefix
	s := NewSerializer()
	value := []byte{0x01, 0x02, 0x03, 0x04}
	s.Bytes(value)
	data := s.ToBytes()

	// Deserialize with zero-copy
	d := NewDeserializer(data)
	got := d.BytesNoCopy()
	if !bytes.Equal(got, value) {
		t.Errorf("BytesNoCopy = %v, want %v", got, value)
	}
	if d.Remaining() != 0 {
		t.Errorf("BytesNoCopy should consume all bytes, remaining: %d", d.Remaining())
	}
}

func TestU128(t *testing.T) {
	tests := []struct {
		name  string
		value *big.Int
	}{
		{"zero", big.NewInt(0)},
		{"one", big.NewInt(1)},
		{"max_u64", new(big.Int).SetUint64(^uint64(0))},
		{"large", func() *big.Int {
			v, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10) // max u128
			return v
		}()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSerializer()
			s.U128(tt.value)
			if s.Error() != nil {
				t.Fatalf("U128 serialize error: %v", s.Error())
			}
			d := NewDeserializer(s.ToBytes())
			got := d.U128()
			if d.Error() != nil {
				t.Fatalf("U128 deserialize error: %v", d.Error())
			}
			if got.Cmp(tt.value) != 0 {
				t.Errorf("U128 roundtrip = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestU256(t *testing.T) {
	tests := []struct {
		name  string
		value *big.Int
	}{
		{"zero", big.NewInt(0)},
		{"one", big.NewInt(1)},
		{"large", func() *big.Int {
			v, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10) // max u256
			return v
		}()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSerializer()
			s.U256(tt.value)
			if s.Error() != nil {
				t.Fatalf("U256 serialize error: %v", s.Error())
			}
			d := NewDeserializer(s.ToBytes())
			got := d.U256()
			if d.Error() != nil {
				t.Fatalf("U256 deserialize error: %v", d.Error())
			}
			if got.Cmp(tt.value) != 0 {
				t.Errorf("U256 roundtrip = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestDeserializerErrors(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		d := NewDeserializer([]byte{})
		_ = d.U8()
		if d.Error() == nil {
			t.Error("expected error for empty input")
		}
	})
	t.Run("invalid bool", func(t *testing.T) {
		d := NewDeserializer([]byte{0x02})
		_ = d.Bool()
		if d.Error() == nil {
			t.Error("expected error for invalid bool")
		}
	})
	t.Run("truncated u64", func(t *testing.T) {
		d := NewDeserializer([]byte{0x01, 0x02, 0x03})
		_ = d.U64()
		if d.Error() == nil {
			t.Error("expected error for truncated u64")
		}
	})
}

// Benchmarks

func BenchmarkSerializerU64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSerializer()
		s.U64(0x123456789abcdef0)
		_ = s.ToBytes()
	}
}

func BenchmarkSerializerU128(b *testing.B) {
	v, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := NewSerializer()
		s.U128(v)
		_ = s.ToBytes()
	}
}

func BenchmarkSerializerString(b *testing.B) {
	str := "hello world"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := NewSerializer()
		s.String(str)
		_ = s.ToBytes()
	}
}

func BenchmarkSerializerUleb128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSerializer()
		s.Uleb128(16384)
		_ = s.ToBytes()
	}
}

func BenchmarkSerializerUleb128Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSerializer()
		s.Uleb128(100)
		_ = s.ToBytes()
	}
}

func BenchmarkDeserializerU64(b *testing.B) {
	data := []byte{0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := NewDeserializer(data)
		_ = d.U64()
	}
}

func BenchmarkDeserializerU128(b *testing.B) {
	// Serialize a U128 value first
	v, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)
	s := NewSerializer()
	s.U128(v)
	data := s.ToBytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := NewDeserializer(data)
		_ = d.U128()
	}
}

func BenchmarkDeserializerU256(b *testing.B) {
	// Serialize a U256 value first
	v, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	s := NewSerializer()
	s.U256(v)
	data := s.ToBytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := NewDeserializer(data)
		_ = d.U256()
	}
}

func BenchmarkAcquireReleaseSerializer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := AcquireSerializer()
		s.U64(0x123456789abcdef0)
		_ = s.ToBytes()
		ReleaseSerializer(s)
	}
}

func BenchmarkAcquireReleaseDeserializer(b *testing.B) {
	data := []byte{0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := AcquireDeserializer(data)
		_ = d.U64()
		ReleaseDeserializer(d)
	}
}

func BenchmarkFixedBytes(b *testing.B) {
	data := make([]byte, 32) // AccountAddress size
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := NewDeserializer(data)
		_ = d.FixedBytes(32)
	}
}

func BenchmarkFixedBytesNoCopy(b *testing.B) {
	data := make([]byte, 32) // AccountAddress size
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := NewDeserializer(data)
		_ = d.FixedBytesNoCopy(32)
	}
}

func BenchmarkBytes(b *testing.B) {
	// Serialize bytes with length prefix
	s := NewSerializer()
	s.Bytes(make([]byte, 64)) // signature size
	data := s.ToBytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := NewDeserializer(data)
		_ = d.Bytes()
	}
}

func BenchmarkBytesNoCopy(b *testing.B) {
	// Serialize bytes with length prefix
	s := NewSerializer()
	s.Bytes(make([]byte, 64)) // signature size
	data := s.ToBytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := NewDeserializer(data)
		_ = d.BytesNoCopy()
	}
}
