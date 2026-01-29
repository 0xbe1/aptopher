package aptos

import (
	"golang.org/x/crypto/sha3"

	"github.com/0xbe1/aptopher/bcs"
	"github.com/0xbe1/aptopher/crypto"
)

// SignedTransaction represents a signed transaction ready for submission.
type SignedTransaction struct {
	RawTxn        *RawTransaction
	Authenticator TransactionAuthenticator
}

// MarshalBCS implements bcs.Marshaler.
func (t SignedTransaction) MarshalBCS(ser *bcs.Serializer) {
	t.RawTxn.MarshalBCS(ser)
	t.Authenticator.MarshalBCS(ser)
}

// Bytes returns the BCS-encoded signed transaction.
func (t *SignedTransaction) Bytes() ([]byte, error) {
	return bcs.Serialize(t)
}

// Hash returns the transaction hash.
// Computes SHA3-256(prefix_hash || variant || signed_txn_bytes) where variant=0 for user transactions.
func (t *SignedTransaction) Hash() (string, error) {
	txnBytes, err := t.Bytes()
	if err != nil {
		return "", err
	}

	// Use incremental hashing to avoid intermediate allocation
	h := sha3.New256()
	h.Write(crypto.TransactionHashPrefix)
	h.Write([]byte{0}) // User transaction variant
	h.Write(txnBytes)

	var hash [32]byte
	h.Sum(hash[:0])
	return bytesToHex(hash[:]), nil
}

func bytesToHex(b []byte) string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, 2+len(b)*2)
	result[0] = '0'
	result[1] = 'x'
	for i, v := range b {
		result[2+i*2] = hexChars[v>>4]
		result[2+i*2+1] = hexChars[v&0x0f]
	}
	return string(result)
}
