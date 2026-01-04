package aptos

import "encoding/json"

// Event represents an on-chain event.
type Event struct {
	GUID           EventGUID       `json:"guid"`
	SequenceNumber string          `json:"sequence_number"`
	Type           string          `json:"type"`
	Data           json.RawMessage `json:"data"`
}

// EventGUID is the globally unique identifier for an event stream.
type EventGUID struct {
	CreationNumber string `json:"creation_number"`
	AccountAddress string `json:"account_address"`
}

// SequenceNumberUint64 returns the sequence number as uint64.
func (e *Event) SequenceNumberUint64() uint64 {
	return parseStringToUint64(e.SequenceNumber)
}

// DecodeData decodes the event data into the provided type.
func (e *Event) DecodeData(v interface{}) error {
	return json.Unmarshal(e.Data, v)
}
