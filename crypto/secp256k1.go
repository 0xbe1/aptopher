package crypto

import (
	"crypto/rand"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

const (
	// Secp256k1PrivateKeyLength is the length of a secp256k1 private key.
	Secp256k1PrivateKeyLength = 32

	// Secp256k1PublicKeyLength is the length of a compressed secp256k1 public key.
	Secp256k1PublicKeyLength = 33

	// Secp256k1SignatureLength is the length of a secp256k1 signature.
	Secp256k1SignatureLength = 64
)

// Secp256k1PrivateKey represents a secp256k1 private key.
type Secp256k1PrivateKey struct {
	key *secp256k1.PrivateKey
}

// GenerateSecp256k1PrivateKey generates a new random secp256k1 private key.
func GenerateSecp256k1PrivateKey() (*Secp256k1PrivateKey, error) {
	keyBytes := make([]byte, Secp256k1PrivateKeyLength)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	key := secp256k1.PrivKeyFromBytes(keyBytes)
	return &Secp256k1PrivateKey{key: key}, nil
}

// NewSecp256k1PrivateKey creates a secp256k1 private key from bytes.
func NewSecp256k1PrivateKey(data []byte) (*Secp256k1PrivateKey, error) {
	if len(data) != Secp256k1PrivateKeyLength {
		return nil, fmt.Errorf("invalid secp256k1 private key length: got %d, want %d", len(data), Secp256k1PrivateKeyLength)
	}
	key := secp256k1.PrivKeyFromBytes(data)
	return &Secp256k1PrivateKey{key: key}, nil
}

// Bytes returns the private key bytes.
func (k *Secp256k1PrivateKey) Bytes() []byte {
	return k.key.Serialize()
}

// Signer returns a Signer for this private key.
func (k *Secp256k1PrivateKey) Signer() Signer {
	return &Secp256k1Signer{key: k.key}
}

// PublicKey returns the compressed public key.
func (k *Secp256k1PrivateKey) PublicKey() []byte {
	return k.key.PubKey().SerializeCompressed()
}

// Secp256k1Signer implements Signer for secp256k1.
type Secp256k1Signer struct {
	key *secp256k1.PrivateKey
}

// Sign signs the message with secp256k1 ECDSA.
// The message is hashed with SHA3-256 before signing.
func (s *Secp256k1Signer) Sign(message []byte) ([]byte, error) {
	// Hash the message with SHA3-256
	hash := Sha3256(message)

	// Sign the hash
	sig := ecdsa.SignCompact(s.key, hash[:], false)
	// SignCompact returns [recovery_id || r || s] (65 bytes)
	// We only need r || s (64 bytes)
	if len(sig) != 65 {
		return nil, fmt.Errorf("unexpected signature length: %d", len(sig))
	}
	return sig[1:], nil // Remove recovery ID
}

// PublicKey returns the compressed secp256k1 public key (33 bytes).
func (s *Secp256k1Signer) PublicKey() []byte {
	return s.key.PubKey().SerializeCompressed()
}

// AuthKey returns the authentication key for this signer.
func (s *Secp256k1Signer) AuthKey() [32]byte {
	return AuthenticationKey(s.PublicKey(), Secp256k1Scheme)
}

// Scheme returns the secp256k1 signature scheme.
func (s *Secp256k1Signer) Scheme() SignatureScheme {
	return Secp256k1Scheme
}

// VerifySecp256k1 verifies a secp256k1 ECDSA signature.
func VerifySecp256k1(publicKey, message, signature []byte) bool {
	if len(publicKey) != Secp256k1PublicKeyLength || len(signature) != Secp256k1SignatureLength {
		return false
	}

	pubKey, err := secp256k1.ParsePubKey(publicKey)
	if err != nil {
		return false
	}

	// Parse signature (r || s format)
	r := new(secp256k1.ModNScalar)
	s := new(secp256k1.ModNScalar)
	r.SetByteSlice(signature[:32])
	s.SetByteSlice(signature[32:])
	sig := ecdsa.NewSignature(r, s)

	// Hash the message
	hash := Sha3256(message)

	return sig.Verify(hash[:], pubKey)
}
