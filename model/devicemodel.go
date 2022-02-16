package model


type DeviceModel struct {
	// Name is DeviceModel name
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Description is DeviceModel description
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Properties is list of DeviceModel properties
	Properties []*Property `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type Property struct {
	// Name is Property name
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// DataType is property dataType
	DataType string `json:"dataType,omitempty" yaml:"dataType,omitempty"`
	// Description is property description
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// AccessMode is property accessMode
	AccessMode string `json:"accessMode,omitempty" yaml:"accessMode,omitempty"`
	// DefaultValue is property defaultValue
	DefaultValue interface{} `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"`
	// Minimum is property minimum value in case of int, double and float
	Minimum interface{} `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	// Maximum is property maximum value in case of int, double and float
	Maximum interface{} `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	// Unit is unit of measurement
	Unit string `json:"unit,omitempty" yaml:"unit,omitempty"`
}