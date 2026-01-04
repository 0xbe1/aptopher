package aptos

import "encoding/json"

// MoveResource represents a Move resource stored on-chain.
type MoveResource struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// DecodeData decodes the resource data into the provided type.
func (r *MoveResource) DecodeData(v interface{}) error {
	return json.Unmarshal(r.Data, v)
}

// CoinStoreResource represents a coin store resource.
type CoinStoreResource struct {
	Coin           CoinResource `json:"coin"`
	DepositEvents  EventHandle  `json:"deposit_events"`
	WithdrawEvents EventHandle  `json:"withdraw_events"`
	Frozen         bool         `json:"frozen"`
}

// CoinResource represents a coin value.
type CoinResource struct {
	Value string `json:"value"`
}

// ValueUint64 returns the coin value as uint64.
func (c *CoinResource) ValueUint64() uint64 {
	return parseStringToUint64(c.Value)
}

// EventHandle represents an event handle.
type EventHandle struct {
	Counter string `json:"counter"`
	GUID    GUID   `json:"guid"`
}

// GUID represents a globally unique identifier for events.
type GUID struct {
	ID GUIDId `json:"id"`
}

// GUIDId contains the creation number and account address.
type GUIDId struct {
	CreationNum string `json:"creation_num"`
	Addr        string `json:"addr"`
}
