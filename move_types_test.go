package aptos

import (
	"encoding/json"
	"testing"
)

func TestU128(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   string
	}{
		{"zero", "0", "0"},
		{"one", "1", "1"},
		{"large", "340282366920938463463374607431768211455", "340282366920938463463374607431768211455"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := U128FromString(tt.input)
			if err != nil {
				t.Fatalf("U128FromString error: %v", err)
			}
			if got := u.String(); got != tt.want {
				t.Errorf("U128.String() = %v, want %v", got, tt.want)
			}

			// JSON roundtrip
			data, err := json.Marshal(u)
			if err != nil {
				t.Fatalf("json.Marshal error: %v", err)
			}
			var got U128
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal error: %v", err)
			}
			if got.String() != tt.want {
				t.Errorf("JSON roundtrip = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestParseTypeTag(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"bool", "bool", "bool", false},
		{"u8", "u8", "u8", false},
		{"u16", "u16", "u16", false},
		{"u32", "u32", "u32", false},
		{"u64", "u64", "u64", false},
		{"u128", "u128", "u128", false},
		{"u256", "u256", "u256", false},
		{"address", "address", "address", false},
		{"signer", "signer", "signer", false},
		{"vector<u8>", "vector<u8>", "vector<u8>", false},
		{"vector<address>", "vector<address>", "vector<address>", false},
		{"simple struct", "0x1::coin::Coin", "0x1::coin::Coin", false},
		{"struct with param", "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>", "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>", false},
		{"nested generics", "0x1::option::Option<0x1::coin::Coin<0x1::aptos_coin::AptosCoin>>", "0x1::option::Option<0x1::coin::Coin<0x1::aptos_coin::AptosCoin>>", false},
		{"invalid", "invalid::type", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTypeTag(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTypeTag(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("ParseTypeTag(%q) = %v, want %v", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestStructTag(t *testing.T) {
	tag := StructTag{
		Address: AccountOne,
		Module:  "coin",
		Name:    "CoinStore",
		TypeParams: []TypeTag{
			{Value: &StructTag{
				Address: AccountOne,
				Module:  "aptos_coin",
				Name:    "AptosCoin",
			}},
		},
	}

	want := "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>"
	if got := tag.String(); got != want {
		t.Errorf("StructTag.String() = %v, want %v", got, want)
	}
}

func TestModuleId(t *testing.T) {
	m := ModuleId{
		Address: AccountOne,
		Name:    "coin",
	}

	want := "0x1::coin"
	if got := m.String(); got != want {
		t.Errorf("ModuleId.String() = %v, want %v", got, want)
	}
}

func TestAptosCoinTypeTag(t *testing.T) {
	want := "0x1::aptos_coin::AptosCoin"
	if got := AptosCoinTypeTag.String(); got != want {
		t.Errorf("AptosCoinTypeTag.String() = %v, want %v", got, want)
	}
}
