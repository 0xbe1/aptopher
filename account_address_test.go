package aptos

import (
	"encoding/json"
	"testing"

	"github.com/0xbe1/lets-go-aptos/bcs"
)

func TestParseAccountAddress(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string // expected full hex string
		wantErr bool
	}{
		{"zero", "0x0", "0x0000000000000000000000000000000000000000000000000000000000000000", false},
		{"one", "0x1", "0x0000000000000000000000000000000000000000000000000000000000000001", false},
		{"short", "0x123", "0x0000000000000000000000000000000000000000000000000000000000000123", false},
		{"full", "0x0000000000000000000000000000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000000000000000000000000000001", false},
		{"no prefix", "1", "0x0000000000000000000000000000000000000000000000000000000000000001", false},
		{"invalid hex", "0xzz", "", true},
		{"too long", "0x" + "ff" + "0000000000000000000000000000000000000000000000000000000000000001", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAccountAddress(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAccountAddress(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("ParseAccountAddress(%q) = %v, want %v", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestAccountAddressShortString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"0x1", "0x1"},
		{"0x0", "0x0"},
		{"0x100", "0x100"},
		{"0x0000000000000000000000000000000000000000000000000000000000000001", "0x1"},
	}
	for _, tt := range tests {
		addr := MustParseAccountAddress(tt.input)
		if got := addr.ShortString(); got != tt.want {
			t.Errorf("AccountAddress(%q).ShortString() = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestAccountAddressJSON(t *testing.T) {
	addr := MustParseAccountAddress("0x1")

	// Marshal
	data, err := json.Marshal(addr)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	// Unmarshal
	var got AccountAddress
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if got != addr {
		t.Errorf("JSON roundtrip: got %v, want %v", got, addr)
	}
}

func TestAccountAddressBCS(t *testing.T) {
	addr := MustParseAccountAddress("0x1")

	// Serialize
	data, err := bcs.Serialize(&addr)
	if err != nil {
		t.Fatalf("BCS serialize error: %v", err)
	}

	// Should be exactly 32 bytes (no length prefix)
	if len(data) != 32 {
		t.Errorf("BCS length = %d, want 32", len(data))
	}

	// Deserialize
	var got AccountAddress
	if err := bcs.Deserialize(data, &got); err != nil {
		t.Fatalf("BCS deserialize error: %v", err)
	}

	if got != addr {
		t.Errorf("BCS roundtrip: got %v, want %v", got, addr)
	}
}

func TestWellKnownAddresses(t *testing.T) {
	if AccountZero.String() != "0x0000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("AccountZero = %v", AccountZero)
	}
	if AccountOne.ShortString() != "0x1" {
		t.Errorf("AccountOne = %v", AccountOne.ShortString())
	}
	if AccountThree.ShortString() != "0x3" {
		t.Errorf("AccountThree = %v", AccountThree.ShortString())
	}
	if AccountFour.ShortString() != "0x4" {
		t.Errorf("AccountFour = %v", AccountFour.ShortString())
	}
}
