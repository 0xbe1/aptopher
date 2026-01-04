package aptos

import (
	"github.com/0xbe1/lets-go-aptos/bcs"
	"github.com/0xbe1/lets-go-aptos/crypto"
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
func (t *SignedTransaction) Hash() (string, error) {
	txnBytes, err := t.Bytes()
	if err != nil {
		return "", err
	}

	// Transaction hash prefix
	prefixHash := crypto.Sha3256([]byte("APTOS::Transaction"))

	// Compute: SHA3-256(prefix_hash || variant || signed_txn_bytes)
	// For user transactions, variant is 0
	data := make([]byte, 0, 33+len(txnBytes))
	data = append(data, prefixHash[:]...)
	data = append(data, 0) // User transaction variant
	data = append(data, txnBytes...)

	hash := crypto.Sha3256(data)
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
