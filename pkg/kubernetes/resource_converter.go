package kubernetes

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

func GetIotPods(ds types.IotDaemonSet, dynamicClient *dynamic.Client,
	restClient *rest.RESTClient) ([]types.IotPod, error) {
	var pods []types.IotPod
	devices, err := GetDaemonSetSelectedDevices(ds, dynamicClient, restClient)

	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		pod := types.IotPod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "IotPod",
				APIVersion: ds.APIVersion,
			},
			Metadata: v1.ObjectMeta{
				Name:      ds.Metadata.Name + "-" + string(common.NewUUID()),
				Namespace: ds.Metadata.Namespace,
				Annotations: map[string]string{
					api.CreatedByAnnotation: ds.Metadata.SelfLink,
					types.DeviceSelector:    device.Metadata.Name,
				},
			},
			Spec: ds.Spec.Template.Spec,
		}
		pods = append(pods, pod)
	}
	return pods, nil
}
