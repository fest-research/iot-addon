package v1

type ResourceKind string

const (
	CreatedBy      = "createdBy"
	DeviceSelector = "deviceSelector"
	DevicesAll     = "all"
	Unschedulable  = "unschedulable"

	APIVersion    = "v1"

	NodeKind     ResourceKind = "Node"
	NodeListKind              = "NodeList"
	PodKind                   = "Pod"
	PodListKind               = "PodList"

	TprIotDevice                 = "iot-device"
	TprIotDaemonSet              = "iot-daemon-set"
	TprIotPod                    = "iot-pod"
)
