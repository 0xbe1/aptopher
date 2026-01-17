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
