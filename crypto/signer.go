// Package crypto provides cryptographic primitives for Aptos transactions.
package crypto

// SignatureScheme represents the signature scheme used.
type SignatureScheme uint8

const (
	// Ed25519Scheme is the Ed25519 signature scheme.
	Ed25519Scheme SignatureScheme = 0

	// Secp256k1Scheme is the secp256k1 ECDSA signature scheme.
	Secp256k1Scheme SignatureScheme = 2
)

// Signer is the interface for signing messages.
type Signer interface {
	// Sign signs the given message and returns the signature.
	Sign(message []byte) ([]byte, error)

	// PublicKey returns the public key bytes.
	PublicKey() []byte

	// AuthKey returns the authentication key derived from the public key.
	AuthKey() [32]byte

	// Scheme returns the signature scheme.
	Scheme() SignatureScheme
}

// PrivateKey represents a private key.
type PrivateKey interface {
	// Bytes returns the private key bytes.
	Bytes() []byte

	// Signer returns a Signer for this private key.
	Signer() Signer
}

// AuthenticationKey derives an authentication key from a public key and scheme.
// For single-key authenticators: SHA3-256(pubkey || scheme)
func AuthenticationKey(pubKey []byte, scheme SignatureScheme) [32]byte {
	data := make([]byte, len(pubKey)+1)
	copy(data, pubKey)
	data[len(pubKey)] = byte(scheme)
	return Sha3256(data)
}
