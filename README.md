### edgex设备属性影子

经过上一阶段对edgex与kubeedge的设备影子调研，edgex支持的所有设备创建设备影子不适合。因为edgex支持的设备属性有数组，对象类型。这些类型在kubeege中目前还没法加入设备影子的属性。并且modbus协议的设备暂时没有条件验证。

因此，作了以下设计前提

（1）、退一步，不以设备为单位创建设备影子。以设备属性为单位创建设备属性影子。（以后，kubeedge设备影子支持数组和对象了，在把该属性赋予创建设备属性影子功能）

（2）、只对mqtt协议的设备属性创建设备属性影子

步骤如下

1、获取device设备的所有devicecommand，每个设备的devicecommand不一定是对设备的全部属性进行操作，只是操作设备属性的其中几个。在设备的所有devicecommand，找出其中操作设备属性全都属于这6种（boolean,bytes,double,float,int,string）中。过滤出设备的devicecommand。

2、针对每个过滤出的devicecommand可以创建一个以该devicecommand名称命名的设备属性影子，设备属性为devicecommand操作的属性。设备属性影子在服务端为k8s crd devicemodel和deviceinstance。

3、（1）在真实设备端，edgex对设备操作后，向设备影子发送特定设备属性的当前真实值。即向kubeedge的mqtt topic $hw/events/device/设备名称/twin/update发送固定格式的属性真实值。

（2）订阅kubeedge的mqtt topic  $hw/events/device/设备名称/twin/update/document，获取特定设备属性的期望值，然后对特定设备属性进行设置，完成后向设备影子发送特定设备属性的当前真实值。即向kubeedge的mqtt topic $hw/events/device/设备名称/twin/update发送固定格式的属性真实值。

下面针对步骤2种，devicecommand种的属性转换为k8s crd devicemodel和deviceinstance构建代码。

devicecomman的数据结构定义如下（edgex的go-mod-core-contracts项目）,可以看到设备属性是个数组，存放类型为ResourceOperation。

```
type DeviceCommand struct {
	Name               string
	IsHidden           bool
	ReadWrite          string
	ResourceOperations []ResourceOperation
}
type ResourceOperation struct {
	DeviceResource string
	DefaultValue   string
	Mappings       map[string]string
}
```

kubeedge devicemodel 数据结构定义如下，需要把数据结构转化为k8s crd的yaml文件。

```
type DeviceModel struct {
	// Name is DeviceModel name
	Name string `json:"name,omitempty"`
	// Description is DeviceModel description
	Description string `json:"description,omitempty"`
	// Properties is list of DeviceModel properties
	Properties []*Property `json:"properties,omitempty"`
}

// Property is structure to store deviceModel property
type Property struct {
	// Name is Property name
	Name string `json:"name,omitempty"`
	// DataType is property dataType
	DataType string `json:"dataType,omitempty"`
	// Description is property description
	Description string `json:"description,omitempty"`
	// AccessMode is property accessMode
	AccessMode string `json:"accessMode,omitempty"`
	// DefaultValue is property defaultValue
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	// Minimum is property minimum value in case of int, double and float
	Minimum interface{} `json:"minimum,omitempty"`
	// Maximum is property maximum value in case of int, double and float
	Maximum interface{} `json:"maximum,omitempty"`
	// Unit is unit of measurement
	Unit string `json:"unit,omitempty"`
}
```

edgex devicecommand元数据与kubeedge devicemodel转换

```
func addDeviceModel(deviceCommand *DeviceCommand, deviceModel *DeviceModel) {
	model := &types.DeviceModel{}
	model.Name = deviceCommand.Name
	model.Properties = make([]*types.Property, 0, len(deviceCommand.ResourceOperations))
	for _, ppt := range deviceCommand.ResourceOperations {
		property := &types.Property{}
		property.Name = ppt.DeviceResource
		property.Description = ppt.DeviceResource
		if ppt.Type.Int != nil {
			property.AccessMode = string(ppt.Type.Int.AccessMode)
			property.DataType = DataTypeInt
			property.DefaultValue = ppt.Type.Int.DefaultValue
			property.Maximum = ppt.Type.Int.Maximum
			property.Minimum = ppt.Type.Int.Minimum
			property.Unit = ppt.Type.Int.Unit
		} else if ppt.Type.String != nil {
			property.AccessMode = string(ppt.Type.String.AccessMode)
			property.DataType = DataTypeString
			property.DefaultValue = ppt.Type.String.DefaultValue
		} else if ppt.Type.Double != nil {
			property.AccessMode = string(ppt.Type.Double.AccessMode)
			property.DataType = DataTypeDouble
			property.DefaultValue = ppt.Type.Double.DefaultValue
			property.Maximum = ppt.Type.Double.Maximum
			property.Minimum = ppt.Type.Double.Minimum
			property.Unit = ppt.Type.Double.Unit
		} else if ppt.Type.Float != nil {
			property.AccessMode = string(ppt.Type.Float.AccessMode)
			property.DataType = DataTypeFloat
			property.DefaultValue = ppt.Type.Float.DefaultValue
			property.Maximum = ppt.Type.Float.Maximum
			property.Minimum = ppt.Type.Float.Minimum
			property.Unit = ppt.Type.Float.Unit
		} else if ppt.Type.Boolean != nil {
			property.AccessMode = string(ppt.Type.Boolean.AccessMode)
			property.DataType = DataTypeBoolean
			property.DefaultValue = ppt.Type.Boolean.DefaultValue
		} else if ppt.Type.Bytes != nil {
			property.AccessMode = string(ppt.Type.Bytes.AccessMode)
			property.DataType = DataTypeBytes
		}
		model.Properties = append(model.Properties, property)
	}
	deviceProfile.DeviceModels = append(deviceProfile.DeviceModels, model)
}
```

