package aptos

import (
	"github.com/0xbe1/lets-go-aptos/bcs"
	"github.com/0xbe1/lets-go-aptos/crypto"
)

// RawTransaction represents an unsigned transaction.
type RawTransaction struct {
	Sender                  AccountAddress
	SequenceNumber          uint64
	Payload                 TransactionPayload
	MaxGasAmount            uint64
	GasUnitPrice            uint64
	ExpirationTimestampSecs uint64
	ChainID                 uint8
}

// MarshalBCS implements bcs.Marshaler.
func (t RawTransaction) MarshalBCS(ser *bcs.Serializer) {
	t.Sender.MarshalBCS(ser)
	ser.U64(t.SequenceNumber)
	t.Payload.MarshalBCS(ser)
	ser.U64(t.MaxGasAmount)
	ser.U64(t.GasUnitPrice)
	ser.U64(t.ExpirationTimestampSecs)
	ser.U8(t.ChainID)
}

// SigningMessage returns the message to be signed for this transaction.
// This is SHA3-256(prefix || bcs(RawTransaction))
func (t *RawTransaction) SigningMessage() ([]byte, error) {
	txnBytes, err := bcs.Serialize(t)
	if err != nil {
		return nil, err
	}
	return crypto.HashWithPrefix(crypto.RawTransactionHashPrefix, txnBytes), nil
}

// Sign signs the transaction with the given signer.
func (t *RawTransaction) Sign(signer crypto.Signer) (*SignedTransaction, error) {
	signingMessage, err := t.SigningMessage()
	if err != nil {
		return nil, err
	}

	signature, err := signer.Sign(signingMessage)
	if err != nil {
		return nil, err
	}

	return &SignedTransaction{
		RawTxn: t,
		Authenticator: TransactionAuthenticator{
			Variant: TransactionAuthenticatorSingleSender,
			Auth: &AccountAuthenticatorSingleKey{
				PublicKey: AnyPublicKey{
					Variant:   signer.Scheme(),
					PublicKey: signer.PublicKey(),
				},
				Signature: AnySignature{
					Variant:   signer.Scheme(),
					Signature: signature,
				},
			},
		},
	}, nil
}

// RawTransactionWithData wraps a raw transaction with additional data for multi-agent/fee-payer transactions.
type RawTransactionWithData struct {
	Variant            RawTransactionWithDataVariant
	RawTxn             *RawTransaction
	SecondarySigners   []AccountAddress // For multi-agent
	FeePayerAddress    AccountAddress   // For fee-payer
}

// RawTransactionWithDataVariant indicates the type of additional data.
type RawTransactionWithDataVariant uint8

const (
	// MultiAgent is for multi-agent transactions.
	MultiAgent RawTransactionWithDataVariant = 0

	// FeePayer is for fee-payer transactions.
	FeePayer RawTransactionWithDataVariant = 1
)

// MarshalBCS implements bcs.Marshaler.
func (t RawTransactionWithData) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(t.Variant))
	t.RawTxn.MarshalBCS(ser)
	switch t.Variant {
	case MultiAgent:
		ser.Uleb128(uint32(len(t.SecondarySigners)))
		for _, addr := range t.SecondarySigners {
			addr.MarshalBCS(ser)
		}
	case FeePayer:
		ser.Uleb128(uint32(len(t.SecondarySigners)))
		for _, addr := range t.SecondarySigners {
			addr.MarshalBCS(ser)
		}
		t.FeePayerAddress.MarshalBCS(ser)
	}
}

// SigningMessage returns the message to be signed for this transaction.
func (t *RawTransactionWithData) SigningMessage() ([]byte, error) {
	txnBytes, err := bcs.Serialize(t)
	if err != nil {
		return nil, err
	}
	return crypto.HashWithPrefix(crypto.RawTransactionWithDataHashPrefix, txnBytes), nil
}
