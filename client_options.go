package aptos

// RequestOptions contains options for API requests.
type RequestOptions struct {
	LedgerVersion *uint64
	Start         *uint64
	Limit         *uint16
}

// RequestOption is a function that modifies request options.
type RequestOption func(*RequestOptions)

// ApplyOptions applies all options and returns the resulting RequestOptions.
func ApplyOptions(opts ...RequestOption) RequestOptions {
	var options RequestOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// WithLedgerVersion specifies a ledger version for the request.
// This retrieves the state at a specific historical version.
func WithLedgerVersion(version uint64) RequestOption {
	return func(o *RequestOptions) {
		o.LedgerVersion = &version
	}
}

// WithStart specifies the starting position for paginated requests.
func WithStart(start uint64) RequestOption {
	return func(o *RequestOptions) {
		o.Start = &start
	}
}

// WithLimit specifies the maximum number of items to return.
func WithLimit(limit uint16) RequestOption {
	return func(o *RequestOptions) {
		o.Limit = &limit
	}
}

// BuildQueryParams builds query parameters from request options.
func (o *RequestOptions) BuildQueryParams() string {
	var params []string
	if o.LedgerVersion != nil {
		params = append(params, "ledger_version="+formatUint64(*o.LedgerVersion))
	}
	if o.Start != nil {
		params = append(params, "start="+formatUint64(*o.Start))
	}
	if o.Limit != nil {
		params = append(params, "limit="+formatUint16(*o.Limit))
	}
	if len(params) == 0 {
		return ""
	}
	return "?" + joinStrings(params, "&")
}

func formatUint64(v uint64) string {
	return formatUint(uint64(v))
}

func formatUint16(v uint16) string {
	return formatUint(uint64(v))
}

func formatUint(v uint64) string {
	// Simple integer to string conversion
	if v == 0 {
		return "0"
	}
	var digits []byte
	for v > 0 {
		digits = append([]byte{byte('0' + v%10)}, digits...)
		v /= 10
	}
	return string(digits)
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}

// BuildOption is a function that modifies transaction build options.
type BuildOption func(*BuildOptions)

// BuildOptions contains options for building transactions.
type BuildOptions struct {
	MaxGasAmount            *uint64
	GasUnitPrice            *uint64
	ExpirationTimestampSecs *uint64
	SequenceNumber          *uint64
	ReplayProtectionNonce   *uint64 // For orderless transactions (mutually exclusive with SequenceNumber)
}

// ApplyBuildOptions applies all build options.
func ApplyBuildOptions(opts ...BuildOption) BuildOptions {
	var options BuildOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// WithMaxGasAmount sets the maximum gas amount for the transaction.
func WithMaxGasAmount(amount uint64) BuildOption {
	return func(o *BuildOptions) {
		o.MaxGasAmount = &amount
	}
}

// WithGasUnitPrice sets the gas unit price for the transaction.
func WithGasUnitPrice(price uint64) BuildOption {
	return func(o *BuildOptions) {
		o.GasUnitPrice = &price
	}
}

// WithExpirationTimestampSecs sets the expiration timestamp for the transaction.
func WithExpirationTimestampSecs(timestamp uint64) BuildOption {
	return func(o *BuildOptions) {
		o.ExpirationTimestampSecs = &timestamp
	}
}

// WithSequenceNumber sets the sequence number for the transaction.
func WithSequenceNumber(seqNum uint64) BuildOption {
	return func(o *BuildOptions) {
		o.SequenceNumber = &seqNum
	}
}

// WithReplayProtectionNonce sets the replay protection nonce for orderless transactions.
// When set, the transaction does not depend on the account's sequence number, allowing
// multiple transactions to be signed and submitted in any order.
// Maximum expiration time for orderless transactions is 60 seconds.
// This option is mutually exclusive with WithSequenceNumber.
func WithReplayProtectionNonce(nonce uint64) BuildOption {
	return func(o *BuildOptions) {
		o.ReplayProtectionNonce = &nonce
	}
}
