package model

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type Device struct {
	DBTimestamp
	Id             string `yaml:"id,omitempty"`
	Name           string `yaml:"name,omitempty"`
	Description    string `yaml:"description,omitempty"`
	AdminState     AdminState `yaml:"adminState,omitempty"`
	OperatingState OperatingState `yaml:"operatingState,omitempty"`
	Protocols      map[string]ProtocolProperties `yaml:"protocols,omitempty"`
	LastConnected  int64 `yaml:"lastConnected,omitempty"`
	LastReported   int64 `yaml:"lastReported,omitempty"`
	Labels         []string `yaml:"labels,omitempty"`
	Location       interface{} `yaml:"location,omitempty"`
	ServiceName    string `yaml:"serviceName,omitempty"`
	ProfileName    string `yaml:"profileName,omitempty"`
	AutoEvents     []AutoEvent `yaml:"autoEvents,omitempty"`
	Notify         bool `yaml:"notify,omitempty"`
}

type ProtocolProperties map[string]string

type AdminState string

type OperatingState string

type AutoEvent struct {
	Interval   string `yaml:"interval,omitempty"`
	OnChange   bool `yaml:"onChange,omitempty"`
	SourceName string `yaml:"sourceName,omitempty"`
}

func (device *Device) GetConf( yamlstr []byte ) *Device {
	err := yaml.UnmarshalStrict(yamlstr,device)
	if err != nil {
		fmt.Println(err.Error())
	}
	return device
}