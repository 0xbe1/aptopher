package aptos

import (
	"encoding/binary"
	"math/big"

	"github.com/0xbe1/aptopher/bcs"
)

// EntryFunctionArg represents a BCS-encoded entry function argument.
type EntryFunctionArg []byte

// EntryFunctionArgs combines multiple arguments into a slice for use in EntryFunction.Args.
func EntryFunctionArgs(args ...EntryFunctionArg) [][]byte {
	result := make([][]byte, len(args))
	for i, arg := range args {
		result[i] = arg
	}
	return result
}

// BoolArg creates a BCS-encoded boolean argument.
func BoolArg(v bool) EntryFunctionArg {
	if v {
		return []byte{1}
	}
	return []byte{0}
}

// U8Arg creates a BCS-encoded u8 argument.
func U8Arg(v uint8) EntryFunctionArg {
	return []byte{v}
}

// U16Arg creates a BCS-encoded u16 argument.
func U16Arg(v uint16) EntryFunctionArg {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, v)
	return buf
}

// U32Arg creates a BCS-encoded u32 argument.
func U32Arg(v uint32) EntryFunctionArg {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, v)
	return buf
}

// U64Arg creates a BCS-encoded u64 argument.
func U64Arg(v uint64) EntryFunctionArg {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, v)
	return buf
}

// U128Arg creates a BCS-encoded u128 argument.
func U128Arg(v *big.Int) EntryFunctionArg {
	ser := bcs.AcquireSerializer()
	ser.U128(v)
	// Must copy since we're releasing the serializer
	result := append([]byte(nil), ser.ToBytes()...)
	bcs.ReleaseSerializer(ser)
	return result
}

// U256Arg creates a BCS-encoded u256 argument.
func U256Arg(v *big.Int) EntryFunctionArg {
	ser := bcs.AcquireSerializer()
	ser.U256(v)
	// Must copy since we're releasing the serializer
	result := append([]byte(nil), ser.ToBytes()...)
	bcs.ReleaseSerializer(ser)
	return result
}

// AddressArg creates a BCS-encoded address argument.
func AddressArg(addr AccountAddress) EntryFunctionArg {
	return addr[:]
}

// StringArg creates a BCS-encoded string argument.
func StringArg(v string) EntryFunctionArg {
	ser := bcs.AcquireSerializer()
	ser.String(v)
	// Must copy since we're releasing the serializer
	result := append([]byte(nil), ser.ToBytes()...)
	bcs.ReleaseSerializer(ser)
	return result
}

// BytesArg creates a BCS-encoded vector<u8> argument.
func BytesArg(v []byte) EntryFunctionArg {
	ser := bcs.AcquireSerializer()
	ser.Bytes(v)
	// Must copy since we're releasing the serializer
	result := append([]byte(nil), ser.ToBytes()...)
	bcs.ReleaseSerializer(ser)
	return result
}

// VectorU8Arg is an alias for BytesArg.
func VectorU8Arg(v []byte) EntryFunctionArg {
	return BytesArg(v)
}

// VectorU64Arg creates a BCS-encoded vector<u64> argument.
func VectorU64Arg(values []uint64) EntryFunctionArg {
	ser := bcs.AcquireSerializer()
	ser.Uleb128(uint32(len(values)))
	for _, v := range values {
		ser.U64(v)
	}
	// Must copy since we're releasing the serializer
	result := append([]byte(nil), ser.ToBytes()...)
	bcs.ReleaseSerializer(ser)
	return result
}

// VectorAddressArg creates a BCS-encoded vector<address> argument.
func VectorAddressArg(addrs []AccountAddress) EntryFunctionArg {
	ser := bcs.AcquireSerializer()
	ser.Uleb128(uint32(len(addrs)))
	for _, addr := range addrs {
		ser.FixedBytes(addr[:])
	}
	// Must copy since we're releasing the serializer
	result := append([]byte(nil), ser.ToBytes()...)
	bcs.ReleaseSerializer(ser)
	return result
}

// VectorStringArg creates a BCS-encoded vector<string> argument.
func VectorStringArg(values []string) EntryFunctionArg {
	ser := bcs.AcquireSerializer()
	ser.Uleb128(uint32(len(values)))
	for _, v := range values {
		ser.String(v)
	}
	// Must copy since we're releasing the serializer
	result := append([]byte(nil), ser.ToBytes()...)
	bcs.ReleaseSerializer(ser)
	return result
}

// OptionU64Arg creates a BCS-encoded Option<u64> argument.
// Pass nil for None, or a pointer to a value for Some.
func OptionU64Arg(v *uint64) EntryFunctionArg {
	ser := bcs.AcquireSerializer()
	if v == nil {
		ser.U8(0) // None
	} else {
		ser.U8(1) // Some
		ser.U64(*v)
	}
	// Must copy since we're releasing the serializer
	result := append([]byte(nil), ser.ToBytes()...)
	bcs.ReleaseSerializer(ser)
	return result
}

// OptionAddressArg creates a BCS-encoded Option<address> argument.
// Pass nil for None, or a pointer to an address for Some.
func OptionAddressArg(addr *AccountAddress) EntryFunctionArg {
	ser := bcs.AcquireSerializer()
	if addr == nil {
		ser.U8(0) // None
	} else {
		ser.U8(1) // Some
		ser.FixedBytes(addr[:])
	}
	// Must copy since we're releasing the serializer
	result := append([]byte(nil), ser.ToBytes()...)
	bcs.ReleaseSerializer(ser)
	return result
}

// OptionStringArg creates a BCS-encoded Option<String> argument.
// Pass nil for None, or a pointer to a string for Some.
func OptionStringArg(v *string) EntryFunctionArg {
	ser := bcs.AcquireSerializer()
	if v == nil {
		ser.U8(0) // None
	} else {
		ser.U8(1) // Some
		ser.String(*v)
	}
	// Must copy since we're releasing the serializer
	result := append([]byte(nil), ser.ToBytes()...)
	bcs.ReleaseSerializer(ser)
	return result
}

// ObjectArg creates a BCS-encoded Object<T> argument (same as address).
func ObjectArg(addr AccountAddress) EntryFunctionArg {
	return AddressArg(addr)
}
