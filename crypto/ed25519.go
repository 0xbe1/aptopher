package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
)

const (
	// Ed25519PublicKeyLength is the length of an Ed25519 public key.
	Ed25519PublicKeyLength = 32

	// Ed25519PrivateKeyLength is the length of an Ed25519 private key (seed).
	Ed25519PrivateKeyLength = 32

	// Ed25519SignatureLength is the length of an Ed25519 signature.
	Ed25519SignatureLength = 64
)

// Ed25519PrivateKey represents an Ed25519 private key.
type Ed25519PrivateKey struct {
	key ed25519.PrivateKey
}

// GenerateEd25519PrivateKey generates a new random Ed25519 private key.
func GenerateEd25519PrivateKey() (*Ed25519PrivateKey, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 key: %w", err)
	}
	return &Ed25519PrivateKey{key: priv}, nil
}

// NewEd25519PrivateKey creates an Ed25519 private key from a 32-byte seed.
func NewEd25519PrivateKey(seed []byte) (*Ed25519PrivateKey, error) {
	if len(seed) != Ed25519PrivateKeyLength {
		return nil, fmt.Errorf("invalid Ed25519 seed length: got %d, want %d", len(seed), Ed25519PrivateKeyLength)
	}
	return &Ed25519PrivateKey{key: ed25519.NewKeyFromSeed(seed)}, nil
}

// Bytes returns the private key seed (32 bytes).
func (k *Ed25519PrivateKey) Bytes() []byte {
	return k.key.Seed()
}

// Signer returns a Signer for this private key.
func (k *Ed25519PrivateKey) Signer() Signer {
	return &Ed25519Signer{key: k.key}
}

// PublicKey returns the public key corresponding to this private key.
func (k *Ed25519PrivateKey) PublicKey() []byte {
	return k.key.Public().(ed25519.PublicKey)
}

// Ed25519Signer implements Signer for Ed25519.
type Ed25519Signer struct {
	key ed25519.PrivateKey
}

// Sign signs the message with Ed25519.
func (s *Ed25519Signer) Sign(message []byte) ([]byte, error) {
	return ed25519.Sign(s.key, message), nil
}

// PublicKey returns the Ed25519 public key (32 bytes).
func (s *Ed25519Signer) PublicKey() []byte {
	return s.key.Public().(ed25519.PublicKey)
}

// AuthKey returns the authentication key for this signer.
func (s *Ed25519Signer) AuthKey() [32]byte {
	return AuthenticationKey(s.PublicKey(), Ed25519Scheme)
}

// Scheme returns the Ed25519 signature scheme.
func (s *Ed25519Signer) Scheme() SignatureScheme {
	return Ed25519Scheme
}

// VerifyEd25519 verifies an Ed25519 signature.
func VerifyEd25519(publicKey, message, signature []byte) bool {
	if len(publicKey) != Ed25519PublicKeyLength || len(signature) != Ed25519SignatureLength {
		return false
	}
	return ed25519.Verify(publicKey, message, signature)
}
