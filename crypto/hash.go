package crypto

import (
	"golang.org/x/crypto/sha3"
)

// Sha3256 computes the SHA3-256 hash of the input.
func Sha3256(data []byte) [32]byte {
	return sha3.Sum256(data)
}

// Sha3256Hash computes the SHA3-256 hash and returns it as a slice.
func Sha3256Hash(data []byte) []byte {
	hash := sha3.Sum256(data)
	return hash[:]
}

// RawTransactionHashPrefix is the prefix used when hashing a raw transaction.
var RawTransactionHashPrefix = sha3256Prefix("APTOS::RawTransaction")

// RawTransactionWithDataHashPrefix is the prefix for transactions with additional data.
var RawTransactionWithDataHashPrefix = sha3256Prefix("APTOS::RawTransactionWithData")

func sha3256Prefix(s string) []byte {
	hash := sha3.Sum256([]byte(s))
	return hash[:]
}

// HashWithPrefix computes SHA3-256(prefix || message).
// Uses incremental hashing to avoid allocating a concatenation buffer.
func HashWithPrefix(prefix, message []byte) []byte {
	h := sha3.New256()
	h.Write(prefix)
	h.Write(message)
	var result [32]byte
	h.Sum(result[:0])
	return result[:]
}
