package aptos

import "encoding/json"

// Block represents an Aptos block.
type Block struct {
	BlockHeight    string            `json:"block_height"`
	BlockHash      string            `json:"block_hash"`
	BlockTimestamp string            `json:"block_timestamp"`
	FirstVersion   string            `json:"first_version"`
	LastVersion    string            `json:"last_version"`
	Transactions   []json.RawMessage `json:"transactions,omitempty"`
}

// BlockHeightUint64 returns the block height as uint64.
func (b *Block) BlockHeightUint64() uint64 {
	return parseStringToUint64(b.BlockHeight)
}

// FirstVersionUint64 returns the first version as uint64.
func (b *Block) FirstVersionUint64() uint64 {
	return parseStringToUint64(b.FirstVersion)
}

// LastVersionUint64 returns the last version as uint64.
func (b *Block) LastVersionUint64() uint64 {
	return parseStringToUint64(b.LastVersion)
}
