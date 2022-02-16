package main

import (
	"edgex-devicetwin/middleware"
	"edgex-devicetwin/model"
	"edgex-devicetwin/v1alpha2"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/yaml.v2"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"

	//v1 "k8s.io/api/core/v1"
	"sort"
)

var INT = []string {common.ValueTypeInt8, common.ValueTypeInt16, common.ValueTypeInt32, common.ValueTypeInt64,
	common.ValueTypeUint8, common.ValueTypeUint16, common.ValueTypeUint32, common.ValueTypeUint64}
var FLOAT = []string {common.ValueTypeFloat32, common.ValueTypeFloat64}
var BOOL = []string{common.ValueTypeBool}
var STRING = []string{common.ValueTypeString}

var ALLTYPE = []string{common.ValueTypeInt8, common.ValueTypeInt16, common.ValueTypeInt32, common.ValueTypeInt64,
	common.ValueTypeUint8, common.ValueTypeUint16, common.ValueTypeUint32, common.ValueTypeUint64,
	common.ValueTypeFloat32, common.ValueTypeFloat64,common.ValueTypeBool,common.ValueTypeString}

const (
	DataTypeInt     = "int"
	DataTypeString  = "string"
	DataTypeFloat   = "float"
	DataTypeBoolean = "boolean"

	Namespace       = "default"
)

type DeviceBody struct {
	Template string `json:"template,omitempty"`
	Data string `json:"data,omitempty"`
}

func in(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}

func availableDeviceCommand(profile *model.DeviceProfile) []model.DeviceCommand {
	var dcs [] model.DeviceCommand
	deviceResources := profile.DeviceResources
	deviceCommands := profile.DeviceCommands
	for _,dc := range deviceCommands {
		pass := true
		passCount := 0
		for _,dcResourceName := range dc.ResourceOperations {
			for _, dr := range deviceResources {
				if dcResourceName.DeviceResource == dr.Name {
					if in(dr.Properties.ValueType,ALLTYPE) {
						passCount ++
					}else {
						pass = false
					}
					break
				}
			}
		}
		if pass && (passCount == len(dc.ResourceOperations)){
			dcs = append(dcs,dc)
		}
	}
	return dcs
}

 func strToFloat64(str string) float64 {
 	r,err := strconv.ParseFloat(str, 32)
 	if err != nil {
 		return 0.0
	}
	return r
 }

 func strToInt(str string) int64 {
	 r,err := strconv.ParseInt(str, 10,32)
	 if err != nil {
		 return 0
	 }
	 return r
 }

 func strToBool(str string) bool {
	 r,err := strconv.ParseBool(str)
	 if err != nil {
		 return false
	 }
	 return r
 }

 func turnToDeviceModel(deviceCommand *model.DeviceCommand,profile *model.DeviceProfile,deviceModel *model.DeviceModel) {
	deviceModel.Name = deviceCommand.Name
	deviceModel.Description= deviceCommand.Name
	deviceModel.Properties = make([]*model.Property, 0, len(deviceCommand.ResourceOperations))
	for _, ppt := range deviceCommand.ResourceOperations {
		property := &model.Property{}
		property.Name = ppt.DeviceResource
		resourcePreperty := &model.ResourceProperties{}
		for _, rp := range profile.DeviceResources{
			if ppt.DeviceResource == rp.Name {
				resourcePreperty.ValueType = rp.Properties.ValueType
				resourcePreperty.ReadWrite = rp.Properties.ReadWrite
				resourcePreperty.Units = rp.Properties.Units
				resourcePreperty.Minimum = rp.Properties.Minimum
				resourcePreperty.Maximum = rp.Properties.Maximum
				resourcePreperty.DefaultValue = rp.Properties.DefaultValue
			}
		}
		property.Description = ppt.DeviceResource
		result := in(resourcePreperty.ValueType, INT)
		if result {
			property.AccessMode = resourcePreperty.ReadWrite
			property.DataType = DataTypeInt
			property.DefaultValue = strToInt(resourcePreperty.DefaultValue)
			property.Maximum = strToInt(resourcePreperty.Maximum)
			property.Minimum = strToInt(resourcePreperty.Minimum)
			property.Unit = resourcePreperty.Units
		}
		result = in(resourcePreperty.ValueType, FLOAT)
		if result {
			property.AccessMode = resourcePreperty.ReadWrite
			property.DataType = DataTypeFloat
			property.DefaultValue = strToFloat64(resourcePreperty.DefaultValue)
			property.Maximum = strToFloat64(resourcePreperty.Maximum)
			property.Minimum = strToFloat64(resourcePreperty.Minimum)
			property.Unit = resourcePreperty.Units
		}
		result = in(resourcePreperty.ValueType, STRING)
		if result {
			property.AccessMode = resourcePreperty.ReadWrite
			property.DataType = DataTypeString
			property.DefaultValue = resourcePreperty.DefaultValue
		}
		result = in(resourcePreperty.ValueType, BOOL)
		if result {
			property.AccessMode = resourcePreperty.ReadWrite
			property.DataType = DataTypeBoolean
			property.DefaultValue = strToBool(resourcePreperty.DefaultValue)
		}
		deviceModel.Properties = append(deviceModel.Properties, property)
	}
}

 func transAccessMode(accessMode string) v1alpha2.PropertyAccessMode {
 	if accessMode == "R" {
 		return v1alpha2.ReadOnly
	}else if accessMode == "RW" {
		return v1alpha2.ReadWrite
	}
	 return v1alpha2.ReadOnly
 }

 func turnToDeviceModelCRD(deviceModel *model.DeviceModel) v1alpha2.DeviceModel {
 	var properties []v1alpha2.DeviceProperty
 	for _,deviceProperties := range deviceModel.Properties {
 		var  propertyType v1alpha2.PropertyType
 		if deviceProperties.DataType == DataTypeFloat {
 			propertyType = v1alpha2.PropertyType{Double: &v1alpha2.PropertyTypeDouble{
				AccessMode: transAccessMode(deviceProperties.AccessMode),
				Maximum:    deviceProperties.Maximum.(float64),
				Minimum: deviceProperties.Minimum.(float64),
				DefaultValue: deviceProperties.DefaultValue.(float64),
				Unit:       deviceProperties.Unit,
			}}
		}else if deviceProperties.DataType == DataTypeInt {
			propertyType = v1alpha2.PropertyType{Int: &v1alpha2.PropertyTypeInt64{
				AccessMode: transAccessMode(deviceProperties.AccessMode),
				Maximum:    deviceProperties.Maximum.(int64),
				Minimum: deviceProperties.Minimum.(int64),
				DefaultValue: deviceProperties.DefaultValue.(int64),
				Unit:       deviceProperties.Unit,
			}}
		}else if deviceProperties.DataType == DataTypeString {
			propertyType = v1alpha2.PropertyType{String: &v1alpha2.PropertyTypeString{
				AccessMode: transAccessMode(deviceProperties.AccessMode),
				DefaultValue: deviceProperties.DefaultValue.(string),
			}}
		}else if deviceProperties.DataType == DataTypeBoolean {
			propertyType = v1alpha2.PropertyType{Boolean: &v1alpha2.PropertyTypeBoolean{
				AccessMode: transAccessMode(deviceProperties.AccessMode),
				DefaultValue: deviceProperties.DefaultValue.(bool),
			}}
		}
 		dp := v1alpha2.DeviceProperty{
		Name:        deviceProperties.Name,
		Description: deviceProperties.Description,
		Type: propertyType,
 		}
 		properties = append(properties,dp)
	}
	 newDeviceModel := v1alpha2.DeviceModel{
		 TypeMeta: v1.TypeMeta{
			 Kind:       "DeviceModel",
			 APIVersion: "devices.kubeedge.io/v1alpha2",
		 },
		 ObjectMeta: v1.ObjectMeta{
			 Name:      deviceModel.Name,
			 Namespace: Namespace,
		 },
		 Spec: v1alpha2.DeviceModelSpec{
			 Properties: properties,
		 },
	 }
	 return newDeviceModel
 }


func turnToDeviceInstance(device *model.Device,deviceModel *model.DeviceModel, nodeSelector string) v1alpha2.Device {

	var twins []v1alpha2.Twin
	twin := v1alpha2.Twin{}
	for _,pr := range deviceModel.Properties {
		twin.PropertyName = pr.Name
		twins = append(twins,twin)
	}

	deviceInstance := v1alpha2.Device{
		TypeMeta: v1.TypeMeta{
			Kind:       "Device",
			APIVersion: "devices.kubeedge.io/v1alpha2",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      device.Name,
			Namespace: Namespace,
		},
		Spec: v1alpha2.DeviceSpec{
			DeviceModelRef: &v12.LocalObjectReference{
				Name: deviceModel.Name,
			},
			NodeSelector: &v12.NodeSelector{
				NodeSelectorTerms: []v12.NodeSelectorTerm{
					{
						MatchExpressions: []v12.NodeSelectorRequirement{
							{
								Key:      "",
								Operator: v12.NodeSelectorOpIn,
								Values:   []string{nodeSelector},
							},
						},
					},
				},
			},
		},
		Status: v1alpha2.DeviceStatus{
			Twins: twins,
		},
	}
	return deviceInstance
}

func main(){
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Any("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Any("/device-translate", devicetranslate)
	_ = r.Run()
}

func devicetranslate(c *gin.Context){
	var deviceBody DeviceBody
	err := c.ShouldBindBodyWith(&deviceBody, binding.JSON)
	if err != nil{
		c.YAML(200, gin.H{"errcode": 400, "description": "Post Data Err"})
		return
	}else {
		fmt.Println(deviceBody.Template)
		fmt.Println(deviceBody.Data)
		var deviceprofile model.DeviceProfile
		var device model.Device
		df := deviceprofile.GetConf([]byte(deviceBody.Template))
		fmt.Println(df)
		de := device.GetConf([]byte(deviceBody.Data))
		fmt.Println(de)
		dcs := availableDeviceCommand(df)
		var dms []model.DeviceModel
		for _, dc :=  range dcs {
		  var dm model.DeviceModel
		  turnToDeviceModel(&dc,df,&dm)
		  dms = append(dms,dm)
		}
		 datastr := "---"
		datastr = datastr + "\n"
		for i, devicemodel := range dms {
			devicemodelcrd := turnToDeviceModelCRD(&devicemodel)
			deviceinstance := turnToDeviceInstance(de,&devicemodel,"edge-test")
			data, _ := yaml.Marshal(devicemodelcrd)
			datastr = datastr + string(data)
			datastr = datastr + "---"
			datastr = datastr + "\n"
			deviceData, _ := yaml.Marshal(deviceinstance)
			datastr = datastr + string(deviceData)
			if i != (len(dms)-1) {
				datastr = datastr + "---"
				datastr = datastr + "\n"
			}
		}
		fmt.Println(datastr)
		c.YAML(200,gin.H{"translate_results": datastr})
	}

}