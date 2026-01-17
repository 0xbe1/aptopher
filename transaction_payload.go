package aptos

import (
	"github.com/0xbe1/aptopher/bcs"
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

	// TransactionPayloadPayload wraps an inner payload (for orderless transactions).
	TransactionPayloadPayload TransactionPayloadVariant = 4
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

// TransactionInnerPayloadV1 wraps an executable with extra config for orderless transactions.
// This is used when replay_protection_nonce is specified instead of sequence_number.
type TransactionInnerPayloadV1 struct {
	Executable  TransactionExecutable
	ExtraConfig TransactionExtraConfigV1
}

func (TransactionInnerPayloadV1) payloadVariant() TransactionPayloadVariant {
	return TransactionPayloadPayload
}

// MarshalBCS implements bcs.Marshaler.
// Serializes as: PayloadVariant(4) + InnerPayloadVariant(0) + Executable + ExtraConfig
func (p TransactionInnerPayloadV1) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(0) // TransactionInnerPayloadVariantV1
	p.Executable.MarshalBCS(ser)
	p.ExtraConfig.MarshalBCS(ser)
}

// TransactionExecutableVariant represents the type of executable.
type TransactionExecutableVariant uint8

const (
	// TransactionExecutableScript is a script executable.
	TransactionExecutableScript TransactionExecutableVariant = 0

	// TransactionExecutableEntryFunction is an entry function executable.
	TransactionExecutableEntryFunction TransactionExecutableVariant = 1
)

// TransactionExecutable wraps a script or entry function for inner payloads.
type TransactionExecutable struct {
	Variant      TransactionExecutableVariant
	Script       *Script
	EntryFunc    *EntryFunction
}

// MarshalBCS implements bcs.Marshaler.
func (e TransactionExecutable) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(e.Variant))
	switch e.Variant {
	case TransactionExecutableScript:
		e.Script.MarshalBCS(ser)
	case TransactionExecutableEntryFunction:
		e.EntryFunc.MarshalBCS(ser)
	}
}

// TransactionExtraConfigV1 contains optional extra configuration for transactions.
type TransactionExtraConfigV1 struct {
	MultisigAddress       *AccountAddress // Optional multisig address
	ReplayProtectionNonce *uint64         // Optional nonce for orderless transactions
}

// MarshalBCS implements bcs.Marshaler.
func (c TransactionExtraConfigV1) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(0) // TransactionExtraConfigVariantV1

	// MultisigAddress as Option<AccountAddress>
	if c.MultisigAddress == nil {
		ser.U8(0) // None
	} else {
		ser.U8(1) // Some
		c.MultisigAddress.MarshalBCS(ser)
	}

	// ReplayProtectionNonce as Option<u64>
	if c.ReplayProtectionNonce == nil {
		ser.U8(0) // None
	} else {
		ser.U8(1) // Some
		ser.U64(*c.ReplayProtectionNonce)
	}
}
