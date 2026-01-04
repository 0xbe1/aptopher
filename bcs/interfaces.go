// Package bcs implements Binary Canonical Serialization (BCS) for Aptos.
// BCS is a non-self-describing binary format used for deterministic serialization.
package bcs

// Marshaler is implemented by types that can serialize themselves to BCS format.
// Unlike encoding/json, this uses a streaming interface to avoid allocations.
type Marshaler interface {
	MarshalBCS(ser *Serializer)
}

// Unmarshaler is implemented by types that can deserialize themselves from BCS format.
type Unmarshaler interface {
	UnmarshalBCS(des *Deserializer)
}

// Struct is implemented by types that support both BCS marshaling and unmarshaling.
type Struct interface {
	Marshaler
	Unmarshaler
}
