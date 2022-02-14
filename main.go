package main

import (
	"edgex-devicetwin/middleware"
	"edgex-devicetwin/model"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/gin-gonic/gin/binding"

	//"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/gin-gonic/gin"
	//"github.com/kubeedge/kubeedge/cloud/pkg/apis/devices/v1alpha2"
	//"github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/types"
	//v1 "k8s.io/api/core/v1"
	"sort"
)

var INT = []string {common.ValueTypeInt8, common.ValueTypeInt16, common.ValueTypeInt32, common.ValueTypeInt64,
	common.ValueTypeUint8, common.ValueTypeUint16, common.ValueTypeUint32, common.ValueTypeUint64}
var FLOAT = []string {common.ValueTypeFloat32, common.ValueTypeFloat64}
var BOOL = []string{common.ValueTypeBool}
var STRING = []string{common.ValueTypeString}

const (
	DataTypeInt     = "int"
	DataTypeString  = "string"
	DataTypeFloat   = "float"
	DataTypeBoolean = "boolean"
)

type DeviceBody struct {
	Template string `json:"template,omitempty"`
	Data string `json:"data,omitempty"`
}

func in(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}


/* func turnToDeviceModel(deviceCommand *models.DeviceCommand,profile *models.DeviceProfile,deviceModel *types.DeviceModel) {
	deviceModel.Name = deviceCommand.Name
	deviceModel.Description= deviceCommand.Name
	deviceModel.Properties = make([]*types.Property, 0, len(deviceCommand.ResourceOperations))
	for _, ppt := range deviceCommand.ResourceOperations {
		property := &types.Property{}
		property.Name = ppt.DeviceResource
		resourcePreperty := &models.ResourceProperties{}
		for _, rp := range profile.DeviceResources{
			if ppt.DeviceResource == rp.Name {
				resourcePreperty = &rp.Properties
			}
		}
		property.Description = ppt.DeviceResource
		result := in(resourcePreperty.ValueType, INT)
		if result {
			property.AccessMode = deviceCommand.ReadWrite
			property.DataType = DataTypeInt
			property.DefaultValue = resourcePreperty.DefaultValue
			property.Maximum = resourcePreperty.Maximum
			property.Minimum = resourcePreperty.Minimum
			property.Unit = resourcePreperty.Units
		}
		result = in(resourcePreperty.ValueType, FLOAT)
		if result {
			property.AccessMode = deviceCommand.ReadWrite
			property.DataType = DataTypeFloat
			property.DefaultValue = resourcePreperty.DefaultValue
			property.Maximum = resourcePreperty.Maximum
			property.Minimum = resourcePreperty.Minimum
			property.Unit = resourcePreperty.Units
		}
		result = in(resourcePreperty.ValueType, STRING)
		if result {
			property.AccessMode = deviceCommand.ReadWrite
			property.DataType = DataTypeString
			property.DefaultValue = resourcePreperty.DefaultValue
		}
		result = in(resourcePreperty.ValueType, BOOL)
		if result {
			property.AccessMode = deviceCommand.ReadWrite
			property.DataType = DataTypeBoolean
			property.DefaultValue = resourcePreperty.DefaultValue
		}
		deviceModel.Properties = append(deviceModel.Properties, property)
	}
}

func turnToDeviceInstance(device *models.Device,deviceService *models.DeviceService,deviceModel *types.DeviceModel,deviceInstance *types.DeviceInstance) {
	deviceInstance.Name = device.Name
	deviceInstance.Model = deviceModel.Name
	nodeIp := deviceService.BaseAddress
	nodeIps := []string{nodeIp}
	nodeSelector := v1.NodeSelector{}
	nodeMatch := v1.NodeSelectorRequirement{Key: "",Operator: "in",Values: nodeIps}
	nodeMatchs := []v1.NodeSelectorRequirement{ nodeMatch}
	nodeSelectorTerm := v1.NodeSelectorTerm{MatchExpressions: nodeMatchs}
	nodeSelector.NodeSelectorTerms = append(nodeSelector.NodeSelectorTerms,nodeSelectorTerm)

	twin := v1alpha2.Twin{}
    for _,pr := range deviceModel.Properties {
    	twin.PropertyName = pr.Name
    	deviceInstance.Twins = append(deviceInstance.Twins,twin)
	}
}*/

func main(){
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Any("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Any("/device-translate", devicetranslate)
	r.Run()
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
		c.YAML(200,gin.H{"render_result": "pong"})
	}

}