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

