package aptos

import (
	"context"
	"fmt"
	"time"
)

// Default transaction parameters
const (
	DefaultMaxGasAmount = uint64(200000)
	DefaultGasUnitPrice = uint64(100)
	DefaultExpirationSeconds = uint64(600) // 10 minutes
)

// BuildTransaction builds a raw transaction for the given sender and payload.
func (c *Client) BuildTransaction(ctx context.Context, sender AccountAddress, payload TransactionPayload, opts ...BuildOption) (*RawTransaction, error) {
	options := ApplyBuildOptions(opts...)

	// Get account info for sequence number
	var sequenceNumber uint64
	if options.SequenceNumber != nil {
		sequenceNumber = *options.SequenceNumber
	} else {
		account, err := c.GetAccount(ctx, sender)
		if err != nil {
			return nil, fmt.Errorf("failed to get account info: %w", err)
		}
		sequenceNumber = account.Data.SequenceNumberUint64()
	}

	// Get gas price if not specified
	var gasUnitPrice uint64
	if options.GasUnitPrice != nil {
		gasUnitPrice = *options.GasUnitPrice
	} else {
		gasEstimate, err := c.EstimateGasPrice(ctx)
		if err != nil {
			// Use default if estimation fails
			gasUnitPrice = DefaultGasUnitPrice
		} else {
			gasUnitPrice = gasEstimate.Data.GasEstimate
		}
	}

	// Get max gas amount
	maxGasAmount := DefaultMaxGasAmount
	if options.MaxGasAmount != nil {
		maxGasAmount = *options.MaxGasAmount
	}

	// Get expiration timestamp
	var expirationTimestampSecs uint64
	if options.ExpirationTimestampSecs != nil {
		expirationTimestampSecs = *options.ExpirationTimestampSecs
	} else {
		expirationTimestampSecs = uint64(time.Now().Unix()) + DefaultExpirationSeconds
	}

	// Get chain ID from ledger info if not cached
	if c.chainID == 0 {
		ledgerInfo, err := c.GetLedgerInfo(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get ledger info: %w", err)
		}
		c.chainID = ledgerInfo.Data.ChainID
	}

	return &RawTransaction{
		Sender:                  sender,
		SequenceNumber:          sequenceNumber,
		Payload:                 payload,
		MaxGasAmount:            maxGasAmount,
		GasUnitPrice:            gasUnitPrice,
		ExpirationTimestampSecs: expirationTimestampSecs,
		ChainID:                 c.chainID,
	}, nil
}

// BuildSignAndSubmitTransaction builds, signs, and submits a transaction.
func (c *Client) BuildSignAndSubmitTransaction(ctx context.Context, account *Account, payload TransactionPayload, opts ...BuildOption) (Response[PendingTransaction], error) {
	// Build transaction
	rawTxn, err := c.BuildTransaction(ctx, account.Address, payload, opts...)
	if err != nil {
		return Response[PendingTransaction]{}, err
	}

	// Sign transaction
	signedTxn, err := account.SignTransaction(rawTxn)
	if err != nil {
		return Response[PendingTransaction]{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Get signed transaction bytes
	txnBytes, err := signedTxn.Bytes()
	if err != nil {
		return Response[PendingTransaction]{}, fmt.Errorf("failed to serialize signed transaction: %w", err)
	}

	// Submit transaction
	return c.SubmitTransaction(ctx, txnBytes)
}

// TransferAPT is a convenience method to transfer APT coins.
func (c *Client) TransferAPT(ctx context.Context, sender *Account, recipient AccountAddress, amount uint64, opts ...BuildOption) (Response[PendingTransaction], error) {
	payload := APTTransferPayload(recipient, amount)
	return c.BuildSignAndSubmitTransaction(ctx, sender, payload, opts...)
}

// SimulateRawTransaction simulates a raw transaction to estimate gas.
func (c *Client) SimulateRawTransaction(ctx context.Context, account *Account, rawTxn *RawTransaction, opts ...SimulateOption) (Response[[]UserTransaction], error) {
	// Create a fake signature for simulation (all zeros)
	fakeSignedTxn := &SignedTransaction{
		RawTxn: rawTxn,
		Authenticator: TransactionAuthenticator{
			Variant: TransactionAuthenticatorSingleSender,
			Auth: &AccountAuthenticatorSingleKey{
				PublicKey: AnyPublicKey{
					Variant:   account.Signer.Scheme(),
					PublicKey: account.Signer.PublicKey(),
				},
				Signature: AnySignature{
					Variant:   account.Signer.Scheme(),
					Signature: make([]byte, 64), // Zero signature for simulation
				},
			},
		},
	}

	txnBytes, err := fakeSignedTxn.Bytes()
	if err != nil {
		return Response[[]UserTransaction]{}, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	return c.SimulateTransaction(ctx, txnBytes, opts...)
}

// WaitForTransaction waits for a transaction to be committed.
// It first tries WaitForTransactionByHash (long-polling), then falls back to polling if that fails.
func (c *Client) WaitForTransaction(ctx context.Context, hash string) (Response[Transaction], error) {
	return c.WaitForTransactionByHash(ctx, hash)
}
