package crypto

import (
	"bytes"
	"testing"
)

func TestEd25519SignAndVerify(t *testing.T) {
	// Generate a key
	priv, err := GenerateEd25519PrivateKey()
	if err != nil {
		t.Fatalf("GenerateEd25519PrivateKey error: %v", err)
	}

	signer := priv.Signer()
	message := []byte("test message")

	// Sign
	sig, err := signer.Sign(message)
	if err != nil {
		t.Fatalf("Sign error: %v", err)
	}

	// Verify
	if !VerifyEd25519(signer.PublicKey(), message, sig) {
		t.Error("signature verification failed")
	}

	// Verify with wrong message should fail
	if VerifyEd25519(signer.PublicKey(), []byte("wrong message"), sig) {
		t.Error("signature verification should have failed")
	}
}

func TestEd25519FromSeed(t *testing.T) {
	seed := make([]byte, Ed25519PrivateKeyLength)
	for i := range seed {
		seed[i] = byte(i)
	}

	priv, err := NewEd25519PrivateKey(seed)
	if err != nil {
		t.Fatalf("NewEd25519PrivateKey error: %v", err)
	}

	// Check seed roundtrip
	if !bytes.Equal(priv.Bytes(), seed) {
		t.Error("seed roundtrip failed")
	}

	// Sign and verify
	signer := priv.Signer()
	message := []byte("test")
	sig, err := signer.Sign(message)
	if err != nil {
		t.Fatalf("Sign error: %v", err)
	}
	if !VerifyEd25519(signer.PublicKey(), message, sig) {
		t.Error("signature verification failed")
	}
}

func TestEd25519AuthKey(t *testing.T) {
	priv, err := GenerateEd25519PrivateKey()
	if err != nil {
		t.Fatalf("GenerateEd25519PrivateKey error: %v", err)
	}

	signer := priv.Signer()
	authKey := signer.AuthKey()

	// Auth key should be 32 bytes
	if len(authKey) != 32 {
		t.Errorf("auth key length = %d, want 32", len(authKey))
	}

	// Auth key should match manual computation
	expected := AuthenticationKey(signer.PublicKey(), Ed25519Scheme)
	if authKey != expected {
		t.Error("auth key mismatch")
	}
}

func TestSecp256k1SignAndVerify(t *testing.T) {
	// Generate a key
	priv, err := GenerateSecp256k1PrivateKey()
	if err != nil {
		t.Fatalf("GenerateSecp256k1PrivateKey error: %v", err)
	}

	signer := priv.Signer()
	message := []byte("test message")

	// Sign
	sig, err := signer.Sign(message)
	if err != nil {
		t.Fatalf("Sign error: %v", err)
	}

	// Check signature length
	if len(sig) != Secp256k1SignatureLength {
		t.Errorf("signature length = %d, want %d", len(sig), Secp256k1SignatureLength)
	}

	// Verify
	if !VerifySecp256k1(signer.PublicKey(), message, sig) {
		t.Error("signature verification failed")
	}

	// Verify with wrong message should fail
	if VerifySecp256k1(signer.PublicKey(), []byte("wrong message"), sig) {
		t.Error("signature verification should have failed")
	}
}

func TestSecp256k1FromBytes(t *testing.T) {
	keyBytes := make([]byte, Secp256k1PrivateKeyLength)
	for i := range keyBytes {
		keyBytes[i] = byte(i + 1)
	}

	priv, err := NewSecp256k1PrivateKey(keyBytes)
	if err != nil {
		t.Fatalf("NewSecp256k1PrivateKey error: %v", err)
	}

	// Public key should be 33 bytes (compressed)
	pubKey := priv.PublicKey()
	if len(pubKey) != Secp256k1PublicKeyLength {
		t.Errorf("public key length = %d, want %d", len(pubKey), Secp256k1PublicKeyLength)
	}

	// Sign and verify
	signer := priv.Signer()
	message := []byte("test")
	sig, err := signer.Sign(message)
	if err != nil {
		t.Fatalf("Sign error: %v", err)
	}
	if !VerifySecp256k1(signer.PublicKey(), message, sig) {
		t.Error("signature verification failed")
	}
}

func TestSha3256(t *testing.T) {
	// Test vector: SHA3-256 of empty string
	hash := Sha3256([]byte{})

	// Known value for SHA3-256("")
	expected := []byte{
		0xa7, 0xff, 0xc6, 0xf8, 0xbf, 0x1e, 0xd7, 0x66,
		0x51, 0xc1, 0x47, 0x56, 0xa0, 0x61, 0xd6, 0x62,
		0xf5, 0x80, 0xff, 0x4d, 0xe4, 0x3b, 0x49, 0xfa,
		0x82, 0xd8, 0x0a, 0x4b, 0x80, 0xf8, 0x43, 0x4a,
	}

	if !bytes.Equal(hash[:], expected) {
		t.Errorf("SHA3-256 hash mismatch")
	}
}
