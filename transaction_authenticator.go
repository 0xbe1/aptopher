package aptos

import (
	"github.com/0xbe1/lets-go-aptos/bcs"
	"github.com/0xbe1/lets-go-aptos/crypto"
)

// TransactionAuthenticatorVariant represents the type of transaction authenticator.
type TransactionAuthenticatorVariant uint8

const (
	// TransactionAuthenticatorEd25519 is a single Ed25519 signature (legacy).
	TransactionAuthenticatorEd25519 TransactionAuthenticatorVariant = 0

	// TransactionAuthenticatorMultiEd25519 is a multi-Ed25519 signature (legacy).
	TransactionAuthenticatorMultiEd25519 TransactionAuthenticatorVariant = 1

	// TransactionAuthenticatorMultiAgent is for multi-agent transactions.
	TransactionAuthenticatorMultiAgent TransactionAuthenticatorVariant = 2

	// TransactionAuthenticatorFeePayer is for fee-payer transactions.
	TransactionAuthenticatorFeePayer TransactionAuthenticatorVariant = 3

	// TransactionAuthenticatorSingleSender is the modern single-sender authenticator.
	TransactionAuthenticatorSingleSender TransactionAuthenticatorVariant = 4
)

// TransactionAuthenticator wraps different authenticator types.
type TransactionAuthenticator struct {
	Variant TransactionAuthenticatorVariant
	Auth    AccountAuthenticatorImpl
}

// MarshalBCS implements bcs.Marshaler.
func (a TransactionAuthenticator) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(a.Variant))
	a.Auth.MarshalBCS(ser)
}

// AccountAuthenticatorImpl is implemented by all account authenticator types.
type AccountAuthenticatorImpl interface {
	bcs.Marshaler
}

// AccountAuthenticatorSingleKey is the modern single-key authenticator.
type AccountAuthenticatorSingleKey struct {
	PublicKey AnyPublicKey
	Signature AnySignature
}

// MarshalBCS implements bcs.Marshaler.
func (a AccountAuthenticatorSingleKey) MarshalBCS(ser *bcs.Serializer) {
	a.PublicKey.MarshalBCS(ser)
	a.Signature.MarshalBCS(ser)
}

// AnyPublicKey represents a public key of any supported type.
type AnyPublicKey struct {
	Variant   crypto.SignatureScheme
	PublicKey []byte
}

// MarshalBCS implements bcs.Marshaler.
func (k AnyPublicKey) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(k.Variant))
	switch k.Variant {
	case crypto.Ed25519Scheme:
		ser.FixedBytes(k.PublicKey) // 32 bytes
	case crypto.Secp256k1Scheme:
		ser.FixedBytes(k.PublicKey) // 33 bytes (compressed)
	default:
		ser.Bytes(k.PublicKey) // Unknown, use length-prefixed
	}
}

// AnySignature represents a signature of any supported type.
type AnySignature struct {
	Variant   crypto.SignatureScheme
	Signature []byte
}

// MarshalBCS implements bcs.Marshaler.
func (s AnySignature) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(s.Variant))
	switch s.Variant {
	case crypto.Ed25519Scheme:
		ser.FixedBytes(s.Signature) // 64 bytes
	case crypto.Secp256k1Scheme:
		ser.FixedBytes(s.Signature) // 64 bytes
	default:
		ser.Bytes(s.Signature) // Unknown, use length-prefixed
	}
}

// AccountAuthenticatorEd25519 is the legacy Ed25519 authenticator.
type AccountAuthenticatorEd25519 struct {
	PublicKey [32]byte
	Signature [64]byte
}

// MarshalBCS implements bcs.Marshaler.
func (a AccountAuthenticatorEd25519) MarshalBCS(ser *bcs.Serializer) {
	ser.FixedBytes(a.PublicKey[:])
	ser.FixedBytes(a.Signature[:])
}

// MultiAgentAuthenticator is for multi-agent transactions.
type MultiAgentAuthenticator struct {
	Sender                   AccountAuthenticatorImpl
	SecondarySignerAddresses []AccountAddress
	SecondarySigners         []AccountAuthenticatorImpl
}

// MarshalBCS implements bcs.Marshaler.
func (a MultiAgentAuthenticator) MarshalBCS(ser *bcs.Serializer) {
	a.Sender.MarshalBCS(ser)
	ser.Uleb128(uint32(len(a.SecondarySignerAddresses)))
	for _, addr := range a.SecondarySignerAddresses {
		addr.MarshalBCS(ser)
	}
	ser.Uleb128(uint32(len(a.SecondarySigners)))
	for _, auth := range a.SecondarySigners {
		auth.MarshalBCS(ser)
	}
}

// FeePayerAuthenticator is for fee-payer transactions.
type FeePayerAuthenticator struct {
	Sender                   AccountAuthenticatorImpl
	SecondarySignerAddresses []AccountAddress
	SecondarySigners         []AccountAuthenticatorImpl
	FeePayerAddress          AccountAddress
	FeePayer                 AccountAuthenticatorImpl
}

// MarshalBCS implements bcs.Marshaler.
func (a FeePayerAuthenticator) MarshalBCS(ser *bcs.Serializer) {
	a.Sender.MarshalBCS(ser)
	ser.Uleb128(uint32(len(a.SecondarySignerAddresses)))
	for _, addr := range a.SecondarySignerAddresses {
		addr.MarshalBCS(ser)
	}
	ser.Uleb128(uint32(len(a.SecondarySigners)))
	for _, auth := range a.SecondarySigners {
		auth.MarshalBCS(ser)
	}
	a.FeePayerAddress.MarshalBCS(ser)
	a.FeePayer.MarshalBCS(ser)
}
