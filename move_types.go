package aptos

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/0xbe1/aptopher/bcs"
)

// U128 represents a 128-bit unsigned integer.
// In JSON, it's represented as a decimal string.
type U128 struct {
	value *big.Int
}

// NewU128 creates a U128 from a uint64.
func NewU128(v uint64) U128 {
	return U128{value: new(big.Int).SetUint64(v)}
}

// NewU128FromBigInt creates a U128 from a big.Int.
func NewU128FromBigInt(v *big.Int) U128 {
	return U128{value: new(big.Int).Set(v)}
}

// U128FromString parses a decimal string into a U128.
func U128FromString(s string) (U128, error) {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return U128{}, fmt.Errorf("invalid U128 string: %s", s)
	}
	return U128{value: v}, nil
}

// BigInt returns the underlying big.Int value.
func (u U128) BigInt() *big.Int {
	if u.value == nil {
		return big.NewInt(0)
	}
	return u.value
}

// Uint64 returns the value as uint64 (may overflow).
func (u U128) Uint64() uint64 {
	if u.value == nil {
		return 0
	}
	return u.value.Uint64()
}

// String returns the decimal string representation.
func (u U128) String() string {
	if u.value == nil {
		return "0"
	}
	return u.value.String()
}

// MarshalJSON implements json.Marshaler.
func (u U128) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *U128) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := U128FromString(s)
	if err != nil {
		return err
	}
	*u = v
	return nil
}

// MarshalBCS implements bcs.Marshaler.
func (u U128) MarshalBCS(ser *bcs.Serializer) {
	ser.U128(u.BigInt())
}

// UnmarshalBCS implements bcs.Unmarshaler.
func (u *U128) UnmarshalBCS(des *bcs.Deserializer) {
	u.value = des.U128()
}

// U256 represents a 256-bit unsigned integer.
type U256 struct {
	value *big.Int
}

// NewU256 creates a U256 from a uint64.
func NewU256(v uint64) U256 {
	return U256{value: new(big.Int).SetUint64(v)}
}

// NewU256FromBigInt creates a U256 from a big.Int.
func NewU256FromBigInt(v *big.Int) U256 {
	return U256{value: new(big.Int).Set(v)}
}

// U256FromString parses a decimal string into a U256.
func U256FromString(s string) (U256, error) {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return U256{}, fmt.Errorf("invalid U256 string: %s", s)
	}
	return U256{value: v}, nil
}

// BigInt returns the underlying big.Int value.
func (u U256) BigInt() *big.Int {
	if u.value == nil {
		return big.NewInt(0)
	}
	return u.value
}

// String returns the decimal string representation.
func (u U256) String() string {
	if u.value == nil {
		return "0"
	}
	return u.value.String()
}

// MarshalJSON implements json.Marshaler.
func (u U256) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *U256) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := U256FromString(s)
	if err != nil {
		return err
	}
	*u = v
	return nil
}

// MarshalBCS implements bcs.Marshaler.
func (u U256) MarshalBCS(ser *bcs.Serializer) {
	ser.U256(u.BigInt())
}

// UnmarshalBCS implements bcs.Unmarshaler.
func (u *U256) UnmarshalBCS(des *bcs.Deserializer) {
	u.value = des.U256()
}

// TypeTagVariant represents the type of a TypeTag.
type TypeTagVariant uint8

const (
	TypeTagBool    TypeTagVariant = 0
	TypeTagU8      TypeTagVariant = 1
	TypeTagU64     TypeTagVariant = 2
	TypeTagU128    TypeTagVariant = 3
	TypeTagAddress TypeTagVariant = 4
	TypeTagSigner  TypeTagVariant = 5
	TypeTagVector  TypeTagVariant = 6
	TypeTagStruct  TypeTagVariant = 7
	TypeTagU16     TypeTagVariant = 8
	TypeTagU32     TypeTagVariant = 9
	TypeTagU256    TypeTagVariant = 10
)

// TypeTag represents a Move type.
type TypeTag struct {
	Value TypeTagValue
}

// TypeTagValue is implemented by all type tag variants.
type TypeTagValue interface {
	bcs.Struct
	typeTagVariant() TypeTagVariant
	String() string
}

// String returns the string representation of the type tag.
func (t TypeTag) String() string {
	if t.Value == nil {
		return ""
	}
	return t.Value.String()
}

// MarshalJSON implements json.Marshaler.
func (t TypeTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *TypeTag) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tag, err := ParseTypeTag(s)
	if err != nil {
		return err
	}
	*t = tag
	return nil
}

// MarshalBCS implements bcs.Marshaler.
func (t TypeTag) MarshalBCS(ser *bcs.Serializer) {
	if t.Value == nil {
		ser.SetError(fmt.Errorf("TypeTag value is nil"))
		return
	}
	ser.Uleb128(uint32(t.Value.typeTagVariant()))
	t.Value.MarshalBCS(ser)
}

// UnmarshalBCS implements bcs.Unmarshaler.
func (t *TypeTag) UnmarshalBCS(des *bcs.Deserializer) {
	variant := TypeTagVariant(des.Uleb128())
	if des.Error() != nil {
		return
	}
	switch variant {
	case TypeTagBool:
		t.Value = &BoolTag{}
	case TypeTagU8:
		t.Value = &U8Tag{}
	case TypeTagU16:
		t.Value = &U16Tag{}
	case TypeTagU32:
		t.Value = &U32Tag{}
	case TypeTagU64:
		t.Value = &U64Tag{}
	case TypeTagU128:
		t.Value = &U128Tag{}
	case TypeTagU256:
		t.Value = &U256Tag{}
	case TypeTagAddress:
		t.Value = &AddressTag{}
	case TypeTagSigner:
		t.Value = &SignerTag{}
	case TypeTagVector:
		var v VectorTag
		v.UnmarshalBCS(des)
		t.Value = &v
	case TypeTagStruct:
		var s StructTag
		s.UnmarshalBCS(des)
		t.Value = &s
	default:
		des.SetError(fmt.Errorf("unknown TypeTag variant: %d", variant))
	}
}

// Primitive type tags
type BoolTag struct{}

func (BoolTag) typeTagVariant() TypeTagVariant   { return TypeTagBool }
func (BoolTag) String() string                   { return "bool" }
func (BoolTag) MarshalBCS(ser *bcs.Serializer)   {}
func (BoolTag) UnmarshalBCS(des *bcs.Deserializer) {}

type U8Tag struct{}

func (U8Tag) typeTagVariant() TypeTagVariant   { return TypeTagU8 }
func (U8Tag) String() string                   { return "u8" }
func (U8Tag) MarshalBCS(ser *bcs.Serializer)   {}
func (U8Tag) UnmarshalBCS(des *bcs.Deserializer) {}

type U16Tag struct{}

func (U16Tag) typeTagVariant() TypeTagVariant   { return TypeTagU16 }
func (U16Tag) String() string                   { return "u16" }
func (U16Tag) MarshalBCS(ser *bcs.Serializer)   {}
func (U16Tag) UnmarshalBCS(des *bcs.Deserializer) {}

type U32Tag struct{}

func (U32Tag) typeTagVariant() TypeTagVariant   { return TypeTagU32 }
func (U32Tag) String() string                   { return "u32" }
func (U32Tag) MarshalBCS(ser *bcs.Serializer)   {}
func (U32Tag) UnmarshalBCS(des *bcs.Deserializer) {}

type U64Tag struct{}

func (U64Tag) typeTagVariant() TypeTagVariant   { return TypeTagU64 }
func (U64Tag) String() string                   { return "u64" }
func (U64Tag) MarshalBCS(ser *bcs.Serializer)   {}
func (U64Tag) UnmarshalBCS(des *bcs.Deserializer) {}

type U128Tag struct{}

func (U128Tag) typeTagVariant() TypeTagVariant   { return TypeTagU128 }
func (U128Tag) String() string                   { return "u128" }
func (U128Tag) MarshalBCS(ser *bcs.Serializer)   {}
func (U128Tag) UnmarshalBCS(des *bcs.Deserializer) {}

type U256Tag struct{}

func (U256Tag) typeTagVariant() TypeTagVariant   { return TypeTagU256 }
func (U256Tag) String() string                   { return "u256" }
func (U256Tag) MarshalBCS(ser *bcs.Serializer)   {}
func (U256Tag) UnmarshalBCS(des *bcs.Deserializer) {}

type AddressTag struct{}

func (AddressTag) typeTagVariant() TypeTagVariant   { return TypeTagAddress }
func (AddressTag) String() string                   { return "address" }
func (AddressTag) MarshalBCS(ser *bcs.Serializer)   {}
func (AddressTag) UnmarshalBCS(des *bcs.Deserializer) {}

type SignerTag struct{}

func (SignerTag) typeTagVariant() TypeTagVariant   { return TypeTagSigner }
func (SignerTag) String() string                   { return "signer" }
func (SignerTag) MarshalBCS(ser *bcs.Serializer)   {}
func (SignerTag) UnmarshalBCS(des *bcs.Deserializer) {}

// VectorTag represents a vector type.
type VectorTag struct {
	ElementType TypeTag
}

func (VectorTag) typeTagVariant() TypeTagVariant { return TypeTagVector }

func (v VectorTag) String() string {
	return fmt.Sprintf("vector<%s>", v.ElementType.String())
}

func (v VectorTag) MarshalBCS(ser *bcs.Serializer) {
	v.ElementType.MarshalBCS(ser)
}

func (v *VectorTag) UnmarshalBCS(des *bcs.Deserializer) {
	v.ElementType.UnmarshalBCS(des)
}

// StructTag represents a struct type.
type StructTag struct {
	Address    AccountAddress
	Module     string
	Name       string
	TypeParams []TypeTag
}

func (StructTag) typeTagVariant() TypeTagVariant { return TypeTagStruct }

func (s StructTag) String() string {
	result := fmt.Sprintf("%s::%s::%s", s.Address.ShortString(), s.Module, s.Name)
	if len(s.TypeParams) > 0 {
		params := make([]string, len(s.TypeParams))
		for i, p := range s.TypeParams {
			params[i] = p.String()
		}
		result += "<" + strings.Join(params, ", ") + ">"
	}
	return result
}

func (s StructTag) MarshalBCS(ser *bcs.Serializer) {
	s.Address.MarshalBCS(ser)
	ser.String(s.Module)
	ser.String(s.Name)
	ser.Uleb128(uint32(len(s.TypeParams)))
	for _, p := range s.TypeParams {
		p.MarshalBCS(ser)
	}
}

func (s *StructTag) UnmarshalBCS(des *bcs.Deserializer) {
	s.Address.UnmarshalBCS(des)
	s.Module = des.String()
	s.Name = des.String()
	length := des.Uleb128()
	s.TypeParams = make([]TypeTag, length)
	for i := uint32(0); i < length; i++ {
		s.TypeParams[i].UnmarshalBCS(des)
	}
}

// ParseTypeTag parses a type tag string into a TypeTag.
// Examples: "u64", "address", "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>"
func ParseTypeTag(s string) (TypeTag, error) {
	s = strings.TrimSpace(s)

	// Primitive types
	switch s {
	case "bool":
		return TypeTag{Value: &BoolTag{}}, nil
	case "u8":
		return TypeTag{Value: &U8Tag{}}, nil
	case "u16":
		return TypeTag{Value: &U16Tag{}}, nil
	case "u32":
		return TypeTag{Value: &U32Tag{}}, nil
	case "u64":
		return TypeTag{Value: &U64Tag{}}, nil
	case "u128":
		return TypeTag{Value: &U128Tag{}}, nil
	case "u256":
		return TypeTag{Value: &U256Tag{}}, nil
	case "address":
		return TypeTag{Value: &AddressTag{}}, nil
	case "signer":
		return TypeTag{Value: &SignerTag{}}, nil
	}

	// Vector type
	if strings.HasPrefix(s, "vector<") && strings.HasSuffix(s, ">") {
		inner := s[7 : len(s)-1]
		elemType, err := ParseTypeTag(inner)
		if err != nil {
			return TypeTag{}, fmt.Errorf("invalid vector element type: %w", err)
		}
		return TypeTag{Value: &VectorTag{ElementType: elemType}}, nil
	}

	// Struct type: address::module::name<type_params>
	return parseStructTag(s)
}

func parseStructTag(s string) (TypeTag, error) {
	// Find type parameters
	var typeParamsStr string
	angleStart := strings.Index(s, "<")
	if angleStart != -1 {
		if !strings.HasSuffix(s, ">") {
			return TypeTag{}, fmt.Errorf("invalid struct tag: missing closing >")
		}
		typeParamsStr = s[angleStart+1 : len(s)-1]
		s = s[:angleStart]
	}

	// Parse address::module::name
	parts := strings.SplitN(s, "::", 3)
	if len(parts) != 3 {
		return TypeTag{}, fmt.Errorf("invalid struct tag format: expected address::module::name")
	}

	addr, err := ParseAccountAddress(parts[0])
	if err != nil {
		return TypeTag{}, fmt.Errorf("invalid struct tag address: %w", err)
	}

	tag := &StructTag{
		Address: addr,
		Module:  parts[1],
		Name:    parts[2],
	}

	// Parse type parameters
	if typeParamsStr != "" {
		params, err := parseTypeParams(typeParamsStr)
		if err != nil {
			return TypeTag{}, err
		}
		tag.TypeParams = params
	}

	return TypeTag{Value: tag}, nil
}

// parseTypeParams parses comma-separated type parameters, handling nested generics.
func parseTypeParams(s string) ([]TypeTag, error) {
	var params []TypeTag
	var depth int
	var start int

	for i, c := range s {
		switch c {
		case '<':
			depth++
		case '>':
			depth--
		case ',':
			if depth == 0 {
				param, err := ParseTypeTag(strings.TrimSpace(s[start:i]))
				if err != nil {
					return nil, err
				}
				params = append(params, param)
				start = i + 1
			}
		}
	}

	// Last parameter
	if start < len(s) {
		param, err := ParseTypeTag(strings.TrimSpace(s[start:]))
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}

	return params, nil
}

// ModuleId identifies a Move module.
type ModuleId struct {
	Address AccountAddress
	Name    string
}

// String returns the module ID as "address::name".
func (m ModuleId) String() string {
	return fmt.Sprintf("%s::%s", m.Address.ShortString(), m.Name)
}

// MarshalBCS implements bcs.Marshaler.
func (m ModuleId) MarshalBCS(ser *bcs.Serializer) {
	m.Address.MarshalBCS(ser)
	ser.String(m.Name)
}

// UnmarshalBCS implements bcs.Unmarshaler.
func (m *ModuleId) UnmarshalBCS(des *bcs.Deserializer) {
	m.Address.UnmarshalBCS(des)
	m.Name = des.String()
}

