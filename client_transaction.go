package aptos

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Default transaction parameters
const (
	DefaultMaxGasAmount      = uint64(200000)
	DefaultGasUnitPrice      = uint64(100)
	DefaultExpirationSeconds = uint64(600) // 10 minutes
)

// BuildTransaction builds a raw transaction for the given sender and payload.
// It fetches required data (sequence number, gas price, chain ID) concurrently
// when not provided via options, reducing latency from 3 round trips to 1.
func (c *Client) BuildTransaction(ctx context.Context, sender AccountAddress, payload TransactionPayload, opts ...BuildOption) (*RawTransaction, error) {
	options := ApplyBuildOptions(opts...)

	// Determine what needs to be fetched
	needSequenceNumber := options.SequenceNumber == nil
	needGasPrice := options.GasUnitPrice == nil
	needChainID := c.chainID == 0

	// Results from concurrent fetches
	var (
		sequenceNumber uint64
		gasUnitPrice   uint64
		chainID        uint8
		fetchErr       error
		mu             sync.Mutex
		wg             sync.WaitGroup
	)

	// Set error helper (first error wins)
	setError := func(err error) {
		mu.Lock()
		if fetchErr == nil {
			fetchErr = err
		}
		mu.Unlock()
	}

	// Fetch sequence number
	if needSequenceNumber {
		wg.Add(1)
		go func() {
			defer wg.Done()
			account, err := c.GetAccount(ctx, sender)
			if err != nil {
				setError(fmt.Errorf("failed to get account info: %w", err))
				return
			}
			mu.Lock()
			sequenceNumber = account.Data.SequenceNumberUint64()
			mu.Unlock()
		}()
	} else {
		sequenceNumber = *options.SequenceNumber
	}

	// Fetch gas price
	if needGasPrice {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gasEstimate, err := c.EstimateGasPrice(ctx)
			mu.Lock()
			if err != nil {
				// Use default if estimation fails (non-fatal)
				gasUnitPrice = DefaultGasUnitPrice
			} else {
				gasUnitPrice = gasEstimate.Data.GasEstimate
			}
			mu.Unlock()
		}()
	} else {
		gasUnitPrice = *options.GasUnitPrice
	}

	// Fetch chain ID
	if needChainID {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ledgerInfo, err := c.GetLedgerInfo(ctx)
			if err != nil {
				setError(fmt.Errorf("failed to get ledger info: %w", err))
				return
			}
			mu.Lock()
			chainID = ledgerInfo.Data.ChainID
			mu.Unlock()
		}()
	} else {
		chainID = c.chainID
	}

	// Wait for all fetches to complete
	wg.Wait()

	// Check for errors
	if fetchErr != nil {
		return nil, fetchErr
	}

	// Cache chain ID for future calls
	if needChainID {
		c.chainID = chainID
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

	return &RawTransaction{
		Sender:                  sender,
		SequenceNumber:          sequenceNumber,
		Payload:                 payload,
		MaxGasAmount:            maxGasAmount,
		GasUnitPrice:            gasUnitPrice,
		ExpirationTimestampSecs: expirationTimestampSecs,
		ChainID:                 chainID,
	}, nil
}

