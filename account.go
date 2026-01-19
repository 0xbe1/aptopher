package aptos

import "strconv"

// AccountData contains basic account information.
type AccountData struct {
	SequenceNumber    string `json:"sequence_number"`
	AuthenticationKey string `json:"authentication_key"`
}

// SequenceNumberUint64 returns the sequence number as uint64.
func (a *AccountData) SequenceNumberUint64() uint64 {
	return parseStringToUint64(a.SequenceNumber)
}

func parseStringToUint64(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}
