package aptos

import "encoding/json"

// MoveModuleBytecode represents a Move module with its bytecode and ABI.
type MoveModuleBytecode struct {
	Bytecode string     `json:"bytecode"`
	ABI      *MoveModule `json:"abi,omitempty"`
}

// MoveModule represents the ABI of a Move module.
type MoveModule struct {
	Address          string           `json:"address"`
	Name             string           `json:"name"`
	Friends          []string         `json:"friends"`
	ExposedFunctions []MoveFunction   `json:"exposed_functions"`
	Structs          []MoveStruct     `json:"structs"`
}

// MoveFunction represents a Move function.
type MoveFunction struct {
	Name              string   `json:"name"`
	Visibility        string   `json:"visibility"`
	IsEntry           bool     `json:"is_entry"`
	IsView            bool     `json:"is_view"`
	GenericTypeParams []MoveFunctionGenericTypeParam `json:"generic_type_params"`
	Params            []string `json:"params"`
	Return            []string `json:"return"`
}

// MoveFunctionGenericTypeParam represents a generic type parameter of a function.
type MoveFunctionGenericTypeParam struct {
	Constraints []string `json:"constraints"`
}

// MoveStruct represents a Move struct.
type MoveStruct struct {
	Name              string   `json:"name"`
	IsNative          bool     `json:"is_native"`
	Abilities         []string `json:"abilities"`
	GenericTypeParams []MoveStructGenericTypeParam `json:"generic_type_params"`
	Fields            []MoveStructField `json:"fields"`
}

// MoveStructGenericTypeParam represents a generic type parameter of a struct.
type MoveStructGenericTypeParam struct {
	Constraints []string                          `json:"constraints"`
	IsPhantom   bool                              `json:"is_phantom"`
}

// MoveStructField represents a field in a Move struct.
type MoveStructField struct {
	Name string          `json:"name"`
	Type json.RawMessage `json:"type"`
}
