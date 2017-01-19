package kubernetes

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func GetIotPod(ds types.IotDaemonSet) (types.IotPod, error) {
	pod := types.IotPod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "IotPod",
			APIVersion: ds.APIVersion,
		},
		Metadata: v1.ObjectMeta{
			Name:      ds.Metadata.Name + "-" + string(common.NewUUID()),
			Namespace: ds.Metadata.Namespace,
		},
	}
	return pod, nil
}
