package model

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type DBTimestamp struct {
	Created  int64 // Created is a timestamp indicating when the entity was created.
	Modified int64 // Modified is a timestamp indicating when the entity was last modified.
}

type DeviceProfile struct {
	DBTimestamp
	Description     string `yaml:"description,omitempty"`
	Id              string `yaml:"id,omitempty"`
	Name            string `yaml:"name,omitempty"`
	Manufacturer    string `yaml:"manufacturer,omitempty"`
	Model           string `yaml:"model,omitempty"`
	Labels          []string `yaml:"labels,omitempty"`
	DeviceResources []DeviceResource `yaml:"deviceResources,omitempty"`
	DeviceCommands  []DeviceCommand `yaml:"deviceCommands,omitempty"`
}

type DeviceResource struct {
	Description string `yaml:"description,omitempty"`
	Name        string `yaml:"name,omitempty"`
	IsHidden    bool `yaml:"isHidden,omitempty"`
	Tag         string `yaml:"tag,omitempty"`
	Properties  ResourceProperties `yaml:"properties,omitempty"`
	Attributes  map[string]interface{} `yaml:"attributes,omitempty"`
}

type ResourceProperties struct {
	ValueType    string `yaml:"valueType,omitempty"`
	ReadWrite    string `yaml:"readWrite,omitempty"`
	Units        string `yaml:"units,omitempty"`
	Minimum      string `yaml:"minimum,omitempty"`
	Maximum      string `yaml:"maximum,omitempty"`
	DefaultValue string `yaml:"defaultValue,omitempty"`
	Mask         string `yaml:"mask,omitempty"`
	Shift        string `yaml:"shift,omitempty"`
	Scale        string `yaml:"scale,omitempty"`
	Offset       string `yaml:"offset,omitempty"`
	Base         string `yaml:"base,omitempty"`
	Assertion    string `yaml:"assertion,omitempty"`
	MediaType    string `yaml:"mediaType,omitempty"`
}

type DeviceCommand struct {
	Name               string `yaml:"name,omitempty"`
	IsHidden           bool `yaml:"isHidden,omitempty"`
	ReadWrite          string `yaml:"readWrite,omitempty"`
	ResourceOperations []ResourceOperation `yaml:"resourceOperations,omitempty"`
}

type ResourceOperation struct {
	DeviceResource string `yaml:"deviceResource,omitempty"`
	DefaultValue   string `yaml:"defaultValue,omitempty"`
	Mappings       map[string]string `yaml:"mappings,omitempty"`
}

func (deviceProfile *DeviceProfile) GetConf( yamlstr []byte ) *DeviceProfile {
	err := yaml.UnmarshalStrict(yamlstr,deviceProfile)
	if err != nil {
		fmt.Println(err.Error())
	}
	return deviceProfile
}
