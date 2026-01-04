package aptos

import (
	"github.com/0xbe1/lets-go-aptos/bcs"
)

// TransactionPayloadVariant represents the type of transaction payload.
type TransactionPayloadVariant uint8

const (
	// TransactionPayloadScript is a script payload.
	TransactionPayloadScript TransactionPayloadVariant = 0

	// TransactionPayloadModuleBundle is deprecated.
	TransactionPayloadModuleBundle TransactionPayloadVariant = 1

	// TransactionPayloadEntryFunction is an entry function payload.
	TransactionPayloadEntryFunction TransactionPayloadVariant = 2

	// TransactionPayloadMultisig is a multisig payload.
	TransactionPayloadMultisig TransactionPayloadVariant = 3
)

// TransactionPayload wraps different payload types.
type TransactionPayload struct {
	Payload TransactionPayloadImpl
}

// TransactionPayloadImpl is implemented by all payload types.
type TransactionPayloadImpl interface {
	bcs.Marshaler
	payloadVariant() TransactionPayloadVariant
}

// MarshalBCS implements bcs.Marshaler.
func (p TransactionPayload) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(p.Payload.payloadVariant()))
	p.Payload.MarshalBCS(ser)
}

// EntryFunction represents an entry function call.
type EntryFunction struct {
	Module   ModuleId
	Function string
	TypeArgs []TypeTag
	Args     [][]byte // BCS-encoded arguments
}

func (EntryFunction) payloadVariant() TransactionPayloadVariant {
	return TransactionPayloadEntryFunction
}

// MarshalBCS implements bcs.Marshaler.
func (e EntryFunction) MarshalBCS(ser *bcs.Serializer) {
	e.Module.MarshalBCS(ser)
	ser.String(e.Function)
	// Type arguments
	ser.Uleb128(uint32(len(e.TypeArgs)))
	for _, t := range e.TypeArgs {
		t.MarshalBCS(ser)
	}
	// Arguments (as vector of bytes)
	ser.Uleb128(uint32(len(e.Args)))
	for _, arg := range e.Args {
		ser.Bytes(arg)
	}
}

// APTTransferPayload creates a payload for transferring APT coins.
func APTTransferPayload(recipient AccountAddress, amount uint64) TransactionPayload {
	return TransactionPayload{
		Payload: &EntryFunction{
			Module: ModuleId{
				Address: AccountOne,
				Name:    "aptos_account",
			},
			Function: "transfer",
			TypeArgs: nil,
			Args: [][]byte{
				mustSerializeAddress(recipient),
				bcs.SerializeU64(amount),
			},
		},
	}
}

// CoinTransferPayload creates a payload for transferring coins of a specific type.
func CoinTransferPayload(coinType TypeTag, recipient AccountAddress, amount uint64) TransactionPayload {
	return TransactionPayload{
		Payload: &EntryFunction{
			Module: ModuleId{
				Address: AccountOne,
				Name:    "aptos_account",
			},
			Function: "transfer_coins",
			TypeArgs: []TypeTag{coinType},
			Args: [][]byte{
				mustSerializeAddress(recipient),
				bcs.SerializeU64(amount),
			},
		},
	}
}

func mustSerializeAddress(addr AccountAddress) []byte {
	data, err := bcs.Serialize(&addr)
	if err != nil {
		panic(err)
	}
	return data
}

// Script represents a Move script.
type Script struct {
	Code     []byte
	TypeArgs []TypeTag
	Args     []ScriptArgument
}

func (Script) payloadVariant() TransactionPayloadVariant {
	return TransactionPayloadScript
}

// MarshalBCS implements bcs.Marshaler.
func (s Script) MarshalBCS(ser *bcs.Serializer) {
	ser.Bytes(s.Code)
	ser.Uleb128(uint32(len(s.TypeArgs)))
	for _, t := range s.TypeArgs {
		t.MarshalBCS(ser)
	}
	ser.Uleb128(uint32(len(s.Args)))
	for _, arg := range s.Args {
		arg.MarshalBCS(ser)
	}
}

// ScriptArgumentVariant represents the type of script argument.
type ScriptArgumentVariant uint8

const (
	ScriptArgumentU8      ScriptArgumentVariant = 0
	ScriptArgumentU64     ScriptArgumentVariant = 1
	ScriptArgumentU128    ScriptArgumentVariant = 2
	ScriptArgumentAddress ScriptArgumentVariant = 3
	ScriptArgumentU8Vec   ScriptArgumentVariant = 4
	ScriptArgumentBool    ScriptArgumentVariant = 5
	ScriptArgumentU16     ScriptArgumentVariant = 6
	ScriptArgumentU32     ScriptArgumentVariant = 7
	ScriptArgumentU256    ScriptArgumentVariant = 8
)

// ScriptArgument represents an argument to a script.
type ScriptArgument struct {
	Variant ScriptArgumentVariant
	Value   interface{}
}

// MarshalBCS implements bcs.Marshaler.
func (a ScriptArgument) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(a.Variant))
	switch a.Variant {
	case ScriptArgumentU8:
		ser.U8(a.Value.(uint8))
	case ScriptArgumentU16:
		ser.U16(a.Value.(uint16))
	case ScriptArgumentU32:
		ser.U32(a.Value.(uint32))
	case ScriptArgumentU64:
		ser.U64(a.Value.(uint64))
	case ScriptArgumentU128:
		v := a.Value.(U128)
		v.MarshalBCS(ser)
	case ScriptArgumentU256:
		v := a.Value.(U256)
		v.MarshalBCS(ser)
	case ScriptArgumentAddress:
		v := a.Value.(AccountAddress)
		v.MarshalBCS(ser)
	case ScriptArgumentU8Vec:
		ser.Bytes(a.Value.([]byte))
	case ScriptArgumentBool:
		ser.Bool(a.Value.(bool))
	}
}

// MultisigPayload represents a multisig transaction payload.
type MultisigPayload struct {
	MultisigAddress    AccountAddress
	TransactionPayload *EntryFunction // Optional
}

func (MultisigPayload) payloadVariant() TransactionPayloadVariant {
	return TransactionPayloadMultisig
}

// MarshalBCS implements bcs.Marshaler.
func (m MultisigPayload) MarshalBCS(ser *bcs.Serializer) {
	m.MultisigAddress.MarshalBCS(ser)
	if m.TransactionPayload == nil {
		ser.U8(0) // None
	} else {
		ser.U8(1) // Some
		m.TransactionPayload.MarshalBCS(ser)
	}
}
