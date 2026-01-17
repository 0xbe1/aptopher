package aptos

// LedgerInfo contains information about the current state of the ledger.
type LedgerInfo struct {
	ChainID             uint8  `json:"chain_id"`
	Epoch               string `json:"epoch"`
	LedgerVersion       string `json:"ledger_version"`
	OldestLedgerVersion string `json:"oldest_ledger_version"`
	LedgerTimestamp     string `json:"ledger_timestamp"`
	NodeRole            string `json:"node_role"`
	OldestBlockHeight   string `json:"oldest_block_height"`
	BlockHeight         string `json:"block_height"`
	GitHash             string `json:"git_hash"`
}

// ResponseMetadata contains metadata from Aptos API response headers.
type ResponseMetadata struct {
	ChainID             uint8
	LedgerVersion       uint64
	LedgerOldestVersion uint64
	LedgerTimestampUsec uint64
	Epoch               uint64
	BlockHeight         uint64
	OldestBlockHeight   uint64
	Cursor              string
}

// Response wraps an API response with metadata from headers.
type Response[T any] struct {
	Data     T
	Metadata ResponseMetadata
}

// BCSResponse wraps raw BCS bytes with response metadata.
type BCSResponse struct {
	Data     []byte
	Metadata ResponseMetadata
}

// NodeInfo contains basic information about a node.
type NodeInfo struct {
	ChainID             uint8  `json:"chain_id"`
	Epoch               string `json:"epoch"`
	LedgerVersion       string `json:"ledger_version"`
	OldestLedgerVersion string `json:"oldest_ledger_version"`
	LedgerTimestamp     string `json:"ledger_timestamp"`
	NodeRole            string `json:"node_role"`
	OldestBlockHeight   string `json:"oldest_block_height"`
	BlockHeight         string `json:"block_height"`
	GitHash             string `json:"git_hash"`
}

// GasEstimation contains gas price estimation from the node.
type GasEstimation struct {
	DeprioritizedGasEstimate uint64 `json:"deprioritized_gas_estimate"`
	GasEstimate              uint64 `json:"gas_estimate"`
	PrioritizedGasEstimate   uint64 `json:"prioritized_gas_estimate"`
}
