package aptos

import "encoding/json"

// Transaction represents an Aptos transaction.
// Use the Type field to determine which specific transaction type this is.
type Transaction struct {
	Type                    string          `json:"type"`
	Hash                    string          `json:"hash"`
	Version                 string          `json:"version,omitempty"`
	StateChangeHash         string          `json:"state_change_hash,omitempty"`
	EventRootHash           string          `json:"event_root_hash,omitempty"`
	StateCheckpointHash     *string         `json:"state_checkpoint_hash,omitempty"`
	GasUsed                 string          `json:"gas_used,omitempty"`
	Success                 bool            `json:"success,omitempty"`
	VMStatus                string          `json:"vm_status,omitempty"`
	AccumulatorRootHash     string          `json:"accumulator_root_hash,omitempty"`
	Changes                 json.RawMessage `json:"changes,omitempty"`
	Sender                  string          `json:"sender,omitempty"`
	SequenceNumber          string          `json:"sequence_number,omitempty"`
	MaxGasAmount            string          `json:"max_gas_amount,omitempty"`
	GasUnitPrice            string          `json:"gas_unit_price,omitempty"`
	ExpirationTimestampSecs string          `json:"expiration_timestamp_secs,omitempty"`
	Payload                 json.RawMessage `json:"payload,omitempty"`
	Signature               json.RawMessage `json:"signature,omitempty"`
	Events                  []Event         `json:"events,omitempty"`
	Timestamp               string          `json:"timestamp,omitempty"`
}

// Transaction types
const (
	TransactionTypePending         = "pending_transaction"
	TransactionTypeUser            = "user_transaction"
	TransactionTypeGenesis         = "genesis_transaction"
	TransactionTypeBlockMetadata   = "block_metadata_transaction"
	TransactionTypeStateCheckpoint = "state_checkpoint_transaction"
	TransactionTypeValidator       = "validator_transaction"
	TransactionTypeBlockEpilogue   = "block_epilogue_transaction"
)

// IsPending returns true if this is a pending transaction.
func (t *Transaction) IsPending() bool {
	return t.Type == TransactionTypePending
}

// IsUserTransaction returns true if this is a user transaction.
func (t *Transaction) IsUserTransaction() bool {
	return t.Type == TransactionTypeUser
}

// VersionUint64 returns the version as uint64.
func (t *Transaction) VersionUint64() uint64 {
	return parseStringToUint64(t.Version)
}

// GasUsedUint64 returns the gas used as uint64.
func (t *Transaction) GasUsedUint64() uint64 {
	return parseStringToUint64(t.GasUsed)
}

// PendingTransaction represents a transaction that has been submitted but not yet committed.
type PendingTransaction struct {
	Hash                    string          `json:"hash"`
	Sender                  string          `json:"sender"`
	SequenceNumber          string          `json:"sequence_number"`
	MaxGasAmount            string          `json:"max_gas_amount"`
	GasUnitPrice            string          `json:"gas_unit_price"`
	ExpirationTimestampSecs string          `json:"expiration_timestamp_secs"`
	Payload                 json.RawMessage `json:"payload"`
	Signature               json.RawMessage `json:"signature"`
}

// UserTransaction represents a committed user transaction.
type UserTransaction struct {
	Version                 string          `json:"version"`
	Hash                    string          `json:"hash"`
	StateChangeHash         string          `json:"state_change_hash"`
	EventRootHash           string          `json:"event_root_hash"`
	StateCheckpointHash     *string         `json:"state_checkpoint_hash"`
	GasUsed                 string          `json:"gas_used"`
	Success                 bool            `json:"success"`
	VMStatus                string          `json:"vm_status"`
	AccumulatorRootHash     string          `json:"accumulator_root_hash"`
	Changes                 json.RawMessage `json:"changes"`
	Sender                  string          `json:"sender"`
	SequenceNumber          string          `json:"sequence_number"`
	MaxGasAmount            string          `json:"max_gas_amount"`
	GasUnitPrice            string          `json:"gas_unit_price"`
	ExpirationTimestampSecs string          `json:"expiration_timestamp_secs"`
	Payload                 json.RawMessage `json:"payload"`
	Signature               json.RawMessage `json:"signature"`
	Events                  []Event         `json:"events"`
	Timestamp               string          `json:"timestamp"`
}

// ViewRequest represents a request to execute a view function.
type ViewRequest struct {
	Function      string        `json:"function"`
	TypeArguments []string      `json:"type_arguments"`
	Arguments     []interface{} `json:"arguments"`
}

// TableItemRequest represents a request to get a table item.
type TableItemRequest struct {
	KeyType   string      `json:"key_type"`
	ValueType string      `json:"value_type"`
	Key       interface{} `json:"key"`
}

// RawTableItemRequest represents a request to get a raw table item.
type RawTableItemRequest struct {
	Key string `json:"key"`
}
