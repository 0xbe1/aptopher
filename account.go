package aptos

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
	var result uint64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + uint64(c-'0')
		}
	}
	return result
}
