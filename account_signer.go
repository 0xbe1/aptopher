package aptos

import (
	"github.com/0xbe1/aptopher/crypto"
)

// Account represents an Aptos account with signing capabilities.
type Account struct {
	Address AccountAddress
	Signer  crypto.Signer
}

// NewEd25519Account generates a new account with a random Ed25519 key.
func NewEd25519Account() (*Account, error) {
	privKey, err := crypto.GenerateEd25519PrivateKey()
	if err != nil {
		return nil, err
	}
	return AccountFromPrivateKey(privKey)
}

// NewSecp256k1Account generates a new account with a random secp256k1 key.
func NewSecp256k1Account() (*Account, error) {
	privKey, err := crypto.GenerateSecp256k1PrivateKey()
	if err != nil {
		return nil, err
	}
	return AccountFromPrivateKey(privKey)
}

// AccountFromPrivateKey creates an account from a private key.
func AccountFromPrivateKey(privKey crypto.PrivateKey) (*Account, error) {
	signer := privKey.Signer()
	authKey := signer.AuthKey()

	var address AccountAddress
	copy(address[:], authKey[:])

	return &Account{
		Address: address,
		Signer:  signer,
	}, nil
}

// AccountFromEd25519Seed creates an account from a 32-byte Ed25519 seed.
func AccountFromEd25519Seed(seed []byte) (*Account, error) {
	privKey, err := crypto.NewEd25519PrivateKey(seed)
	if err != nil {
		return nil, err
	}
	return AccountFromPrivateKey(privKey)
}

// AccountFromSecp256k1Bytes creates an account from a 32-byte secp256k1 private key.
func AccountFromSecp256k1Bytes(keyBytes []byte) (*Account, error) {
	privKey, err := crypto.NewSecp256k1PrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return AccountFromPrivateKey(privKey)
}

// Sign signs a message with this account's private key.
func (a *Account) Sign(message []byte) ([]byte, error) {
	return a.Signer.Sign(message)
}

// SignTransaction signs a raw transaction.
func (a *Account) SignTransaction(rawTxn *RawTransaction) (*SignedTransaction, error) {
	return rawTxn.Sign(a.Signer)
}

// AuthKey returns the authentication key for this account.
func (a *Account) AuthKey() [32]byte {
	return a.Signer.AuthKey()
}
