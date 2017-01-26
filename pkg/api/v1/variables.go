package v1

type ResourceKind string

const (
	CreatedBy      = "createdBy"
	DeviceSelector = "deviceSelector"
	DevicesAll     = "all"
	Unschedulable  = "unschedulable"

	IotAPIVersion = "fujitsu.com/v1"
	APIVersion    = "v1"

	NodeKind     ResourceKind = "Node"
	NodeListKind              = "NodeList"
	PodKind                   = "Pod"
	PodListKind               = "PodList"
)
