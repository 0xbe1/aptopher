package aptos

import (
	"encoding/json"
	"fmt"

	"github.com/0xbe1/aptopher/bcs"
	"github.com/0xbe1/aptopher/internal/hex"
)

const (
	// AccountAddressLength is the length of an Aptos account address in bytes.
	AccountAddressLength = 32
)

// AccountAddress represents a 32-byte Aptos account address.
type AccountAddress [AccountAddressLength]byte

// Well-known addresses
var (
	// AccountZero is the zero address (0x0).
	AccountZero = AccountAddress{}

	// AccountOne is the core framework address (0x1).
	AccountOne = mustParseAddress("0x1")

	// AccountThree is the token address (0x3).
	AccountThree = mustParseAddress("0x3")

	// AccountFour is the token objects address (0x4).
	AccountFour = mustParseAddress("0x4")
)

func mustParseAddress(s string) AccountAddress {
	addr, err := ParseAccountAddress(s)
	if err != nil {
		panic(err)
	}
	return addr
}

// ParseAccountAddress parses a hex string (with or without 0x prefix) into an AccountAddress.
// Short addresses are left-padded with zeros.
func ParseAccountAddress(s string) (AccountAddress, error) {
	data, err := hex.DecodeFixed(s, AccountAddressLength)
	if err != nil {
		return AccountAddress{}, fmt.Errorf("invalid account address %q: %w", s, err)
	}
	var addr AccountAddress
	copy(addr[:], data)
	return addr, nil
}

// MustParseAccountAddress parses an address or panics.
func MustParseAccountAddress(s string) AccountAddress {
	addr, err := ParseAccountAddress(s)
	if err != nil {
		panic(err)
	}
	return addr
}

// String returns the address as a 0x-prefixed hex string.
// Leading zeros are preserved for the full 32-byte representation.
func (a AccountAddress) String() string {
	return hex.Encode(a[:])
}

// ShortString returns a shortened address string with leading zeros removed.
// Example: "0x1" instead of "0x0000...0001"
func (a AccountAddress) ShortString() string {
	s := hex.Encode(a[:])
	// Strip leading zeros after "0x", but keep at least one digit
	s = s[2:] // Remove "0x" prefix
	for len(s) > 1 && s[0] == '0' {
		s = s[1:]
	}
	return "0x" + s
}

// Bytes returns the address as a byte slice.
func (a AccountAddress) Bytes() []byte {
	return a[:]
}

// IsZero returns true if this is the zero address.
func (a AccountAddress) IsZero() bool {
	return a == AccountZero
}

// MarshalJSON implements json.Marshaler.
func (a AccountAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AccountAddress) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	addr, err := ParseAccountAddress(s)
	if err != nil {
		return err
	}
	*a = addr
	return nil
}

// MarshalBCS implements bcs.Marshaler.
// AccountAddress is serialized as a fixed 32-byte array (no length prefix).
func (a AccountAddress) MarshalBCS(ser *bcs.Serializer) {
	ser.FixedBytes(a[:])
}

// UnmarshalBCS implements bcs.Unmarshaler.
func (a *AccountAddress) UnmarshalBCS(des *bcs.Deserializer) {
	data := des.FixedBytes(AccountAddressLength)
	if des.Error() != nil {
		return
	}
	copy(a[:], data)
}
