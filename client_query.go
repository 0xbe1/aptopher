package aptos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// GetTransactions retrieves a list of transactions.
func (c *Client) GetTransactions(ctx context.Context, opts ...RequestOption) (Response[[]Transaction], error) {
	options := ApplyOptions(opts...)
	path := "/transactions" + options.BuildQueryParams()

	var txns []Transaction
	metadata, err := c.http.get(ctx, path, &txns)
	if err != nil {
		return Response[[]Transaction]{}, err
	}
	return Response[[]Transaction]{Data: txns, Metadata: metadata}, nil
}

// GetTransactionByHash retrieves a transaction by its hash.
func (c *Client) GetTransactionByHash(ctx context.Context, hash string) (Response[Transaction], error) {
	path := "/transactions/by_hash/" + hash

	var txn Transaction
	metadata, err := c.http.get(ctx, path, &txn)
	if err != nil {
		return Response[Transaction]{}, err
	}
	return Response[Transaction]{Data: txn, Metadata: metadata}, nil
}

// WaitForTransactionByHash waits for a transaction to be committed.
// This uses long-polling and will block until the transaction is committed or times out.
func (c *Client) WaitForTransactionByHash(ctx context.Context, hash string) (Response[Transaction], error) {
	path := "/transactions/wait_by_hash/" + hash

	var txn Transaction
	metadata, err := c.http.get(ctx, path, &txn)
	if err != nil {
		return Response[Transaction]{}, err
	}
	return Response[Transaction]{Data: txn, Metadata: metadata}, nil
}

// GetTransactionByVersion retrieves a transaction by its version.
func (c *Client) GetTransactionByVersion(ctx context.Context, version uint64) (Response[Transaction], error) {
	path := fmt.Sprintf("/transactions/by_version/%d", version)

	var txn Transaction
	metadata, err := c.http.get(ctx, path, &txn)
	if err != nil {
		return Response[Transaction]{}, err
	}
	return Response[Transaction]{Data: txn, Metadata: metadata}, nil
}

// GetAccountTransactions retrieves transactions for a specific account.
func (c *Client) GetAccountTransactions(ctx context.Context, address AccountAddress, opts ...RequestOption) (Response[[]Transaction], error) {
	options := ApplyOptions(opts...)
	path := "/accounts/" + address.String() + "/transactions" + options.BuildQueryParams()

	var txns []Transaction
	metadata, err := c.http.get(ctx, path, &txns)
	if err != nil {
		return Response[[]Transaction]{}, err
	}
	return Response[[]Transaction]{Data: txns, Metadata: metadata}, nil
}

// GetBlockByHeight retrieves a block by its height.
func (c *Client) GetBlockByHeight(ctx context.Context, height uint64, withTransactions bool) (Response[Block], error) {
	path := fmt.Sprintf("/blocks/by_height/%d", height)
	if withTransactions {
		path += "?with_transactions=true"
	}

	var block Block
	metadata, err := c.http.get(ctx, path, &block)
	if err != nil {
		return Response[Block]{}, err
	}
	return Response[Block]{Data: block, Metadata: metadata}, nil
}

// GetBlockByVersion retrieves a block by the ledger version it contains.
func (c *Client) GetBlockByVersion(ctx context.Context, version uint64, withTransactions bool) (Response[Block], error) {
	path := fmt.Sprintf("/blocks/by_version/%d", version)
	if withTransactions {
		path += "?with_transactions=true"
	}

	var block Block
	metadata, err := c.http.get(ctx, path, &block)
	if err != nil {
		return Response[Block]{}, err
	}
	return Response[Block]{Data: block, Metadata: metadata}, nil
}

// GetEventsByCreationNumber retrieves events by creation number.
func (c *Client) GetEventsByCreationNumber(ctx context.Context, address AccountAddress, creationNumber uint64, opts ...RequestOption) (Response[[]Event], error) {
	options := ApplyOptions(opts...)
	path := fmt.Sprintf("/accounts/%s/events/%d%s", address.String(), creationNumber, options.BuildQueryParams())

	var events []Event
	metadata, err := c.http.get(ctx, path, &events)
	if err != nil {
		return Response[[]Event]{}, err
	}
	return Response[[]Event]{Data: events, Metadata: metadata}, nil
}

// GetEventsByEventHandle retrieves events by event handle.
func (c *Client) GetEventsByEventHandle(ctx context.Context, address AccountAddress, eventHandle, fieldName string, opts ...RequestOption) (Response[[]Event], error) {
	options := ApplyOptions(opts...)
	path := fmt.Sprintf("/accounts/%s/events/%s/%s%s",
		address.String(),
		url.PathEscape(eventHandle),
		url.PathEscape(fieldName),
		options.BuildQueryParams())

	var events []Event
	metadata, err := c.http.get(ctx, path, &events)
	if err != nil {
		return Response[[]Event]{}, err
	}
	return Response[[]Event]{Data: events, Metadata: metadata}, nil
}

// GetTableItem retrieves a table item.
func (c *Client) GetTableItem(ctx context.Context, tableHandle string, req TableItemRequest, opts ...RequestOption) (Response[json.RawMessage], error) {
	options := ApplyOptions(opts...)
	path := "/tables/" + tableHandle + "/item" + options.BuildQueryParams()

	var result json.RawMessage
	metadata, err := c.http.post(ctx, path, req, &result)
	if err != nil {
		return Response[json.RawMessage]{}, err
	}
	return Response[json.RawMessage]{Data: result, Metadata: metadata}, nil
}

// GetRawTableItem retrieves a raw table item.
func (c *Client) GetRawTableItem(ctx context.Context, tableHandle string, req RawTableItemRequest, opts ...RequestOption) (Response[json.RawMessage], error) {
	options := ApplyOptions(opts...)
	path := "/tables/" + tableHandle + "/raw_item" + options.BuildQueryParams()

	var result json.RawMessage
	metadata, err := c.http.post(ctx, path, req, &result)
	if err != nil {
		return Response[json.RawMessage]{}, err
	}
	return Response[json.RawMessage]{Data: result, Metadata: metadata}, nil
}

// View executes a view function and returns the result.
func (c *Client) View(ctx context.Context, req ViewRequest, opts ...RequestOption) (Response[[]json.RawMessage], error) {
	options := ApplyOptions(opts...)
	path := "/view" + options.BuildQueryParams()

	var result []json.RawMessage
	metadata, err := c.http.post(ctx, path, req, &result)
	if err != nil {
		return Response[[]json.RawMessage]{}, err
	}
	return Response[[]json.RawMessage]{Data: result, Metadata: metadata}, nil
}

// SimulateTransaction simulates a transaction without committing it.
func (c *Client) SimulateTransaction(ctx context.Context, signedTxnBytes []byte, opts ...SimulateOption) (Response[[]UserTransaction], error) {
	simOpts := ApplySimulateOptions(opts...)
	path := "/transactions/simulate"

	var params []string
	if simOpts.EstimateMaxGasAmount {
		params = append(params, "estimate_max_gas_amount=true")
	}
	if simOpts.EstimateGasUnitPrice {
		params = append(params, "estimate_gas_unit_price=true")
	}
	if simOpts.EstimatePrioritizedGasUnitPrice {
		params = append(params, "estimate_prioritized_gas_unit_price=true")
	}
	if len(params) > 0 {
		path += "?" + joinStrings(params, "&")
	}

	var result []UserTransaction
	metadata, err := c.http.postBCS(ctx, path, signedTxnBytes, &result)
	if err != nil {
		return Response[[]UserTransaction]{}, err
	}
	return Response[[]UserTransaction]{Data: result, Metadata: metadata}, nil
}

// SubmitTransaction submits a signed transaction.
func (c *Client) SubmitTransaction(ctx context.Context, signedTxnBytes []byte) (Response[PendingTransaction], error) {
	path := "/transactions"

	var result PendingTransaction
	metadata, err := c.http.postBCS(ctx, path, signedTxnBytes, &result)
	if err != nil {
		return Response[PendingTransaction]{}, err
	}
	return Response[PendingTransaction]{Data: result, Metadata: metadata}, nil
}


// SimulateOption is a function that modifies simulation options.
type SimulateOption func(*SimulateOptions)

// SimulateOptions contains options for transaction simulation.
type SimulateOptions struct {
	EstimateMaxGasAmount           bool
	EstimateGasUnitPrice           bool
	EstimatePrioritizedGasUnitPrice bool
}

// ApplySimulateOptions applies all simulation options.
func ApplySimulateOptions(opts ...SimulateOption) SimulateOptions {
	var options SimulateOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// WithEstimateMaxGasAmount enables max gas amount estimation.
func WithEstimateMaxGasAmount() SimulateOption {
	return func(o *SimulateOptions) {
		o.EstimateMaxGasAmount = true
	}
}

// WithEstimateGasUnitPrice enables gas unit price estimation.
func WithEstimateGasUnitPrice() SimulateOption {
	return func(o *SimulateOptions) {
		o.EstimateGasUnitPrice = true
	}
}

// WithEstimatePrioritizedGasUnitPrice enables prioritized gas unit price estimation.
func WithEstimatePrioritizedGasUnitPrice() SimulateOption {
	return func(o *SimulateOptions) {
		o.EstimatePrioritizedGasUnitPrice = true
	}
}

// PollForTransaction polls for a transaction until it's found or the context is cancelled.
// This is useful when long-polling is not available or times out.
func (c *Client) PollForTransaction(ctx context.Context, hash string, pollInterval time.Duration) (Response[Transaction], error) {
	for {
		txn, err := c.GetTransactionByHash(ctx, hash)
		if err == nil && !txn.Data.IsPending() {
			return txn, nil
		}

		select {
		case <-ctx.Done():
			return Response[Transaction]{}, ctx.Err()
		case <-time.After(pollInterval):
			// Continue polling
		}
	}
}
