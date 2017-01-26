package controller

import (
	"reflect"
	"testing"

	"github.com/fest-research/iot-addon/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	kubeapi "k8s.io/client-go/pkg/api/v1"
)

var (
	ctrl = NewNodeController()

	iotTypeMeta = metav1.TypeMeta{
		Kind:       string(v1.IotDeviceKind),
		APIVersion: v1.IotAPIVersion,
	}

	typeMeta = metav1.TypeMeta{
		Kind:       string(v1.NodeKind),
		APIVersion: v1.APIVersion,
	}
)

func createTestNode(name string) *kubeapi.Node {
	return &kubeapi.Node{
		TypeMeta: typeMeta,
		ObjectMeta: kubeapi.ObjectMeta{
			Name:      name,
			Namespace: "",
		},
	}
}

func createTestIotDevice(name, namespace string) *v1.IotDevice {
	return &v1.IotDevice{
		TypeMeta: iotTypeMeta,
		Metadata: kubeapi.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

// TODO: test more cases (i.e. null objects handling)

func TestTransformWatchEvent(t *testing.T) {
	node := createTestNode("test-node")
	iotDevice := createTestIotDevice("test-node", "default")

	nodeObj := runtime.Object(node)
	iotDeviceObj := runtime.Object(iotDevice)

	cases := []struct {
		event    watch.Event
		expected watch.Event
	}{
		{
			event:    watch.Event{Object: iotDeviceObj},
			expected: watch.Event{Object: nodeObj},
		},
	}

	for _, c := range cases {
		result := ctrl.TransformWatchEvent(c.event)

		if reflect.TypeOf(result.Object) != reflect.TypeOf(c.expected.Object) {
			t.Errorf("TransformWatchEvent(event: %v): expected: %s, got: %s", c.event,
				reflect.TypeOf(c.expected.Object), reflect.TypeOf(result.Object))
		}
	}
}

func TestToNodeList(t *testing.T) {
	iotDevice := createTestIotDevice("test-node", "default-ns")
	node := createTestNode("test-node")

	cases := []struct {
		iotDeviceList    *v1.IotDeviceList
		expected         *kubeapi.NodeList
		expectedTypeMeta metav1.TypeMeta
	}{
		{
			&v1.IotDeviceList{Items: []v1.IotDevice{*iotDevice}},
			&kubeapi.NodeList{Items: []kubeapi.Node{*node}},
			typeMeta,
		},
	}

	for _, c := range cases {
		result := ctrl.ToNodeList(c.iotDeviceList)

		if reflect.TypeOf(c.expected) != reflect.TypeOf(result) {
			t.Errorf("ToNodeList(iotDeviceList: %v): expected: %s, got: %s", c.iotDeviceList,
				reflect.TypeOf(c.expected), reflect.TypeOf(result))
		}

		if len(c.expected.Items) != len(result.Items) {
			t.Errorf("ToNodeList(iotDeviceList: %v): expected: %d, got: %d", c.iotDeviceList,
				len(c.expected.Items), len(result.Items))
		}

		for _, node := range result.Items {
			if !reflect.DeepEqual(node.TypeMeta, c.expectedTypeMeta) {
				t.Errorf("ToNodeList(iotDeviceList: %v): expected: %v, got: %v", c.iotDeviceList,
					c.expectedTypeMeta, node.TypeMeta)
			}
		}
	}
}

func TestToNode(t *testing.T) {
	cases := []struct {
		iotDevice        *v1.IotDevice
		expected         *kubeapi.Node
		expectedTypeMeta metav1.TypeMeta
	}{
		{
			createTestIotDevice("test-node", "default-ns"),
			createTestNode("test-node"),
			typeMeta,
		},
	}

	for _, c := range cases {
		result := ctrl.ToNode(c.iotDevice)

		if !reflect.DeepEqual(result.TypeMeta, c.expectedTypeMeta) {
			t.Errorf("ToNode(iotDevice: %v): expected: %v, got: %v", c.iotDevice,
				c.expectedTypeMeta, result.TypeMeta)
		}

		if !reflect.DeepEqual(result.Spec, c.expected.Spec) ||
			!reflect.DeepEqual(result.Status, c.expected.Status) ||
			!reflect.DeepEqual(result.ObjectMeta, c.expected.ObjectMeta) ||
			result.Namespace != "" {
			t.Errorf("ToNode(iotDevice: %v): expected: %v, got: %v", c.iotDevice,
				c.expected, result)
		}
	}
}

func TestToIotDevice(t *testing.T) {
	cases := []struct {
		node             *kubeapi.Node
		expected         *v1.IotDevice
		expectedTypeMeta metav1.TypeMeta
	}{
		{
			createTestNode("test-node"),
			createTestIotDevice("test-node", "default-ns"),
			iotTypeMeta,
		},
	}

	for _, c := range cases {
		result := ctrl.ToIotDevice(c.node)

		if !reflect.DeepEqual(result.TypeMeta, c.expectedTypeMeta) {
			t.Errorf("ToIotDevice(node: %v): expected: %v, got: %v", c.node,
				c.expectedTypeMeta, result.TypeMeta)
		}

		// TODO: test namespace?
		if !reflect.DeepEqual(result.Spec, c.expected.Spec) ||
			!reflect.DeepEqual(result.Status, c.expected.Status) {
			t.Errorf("ToIotDevice(node: %v): expected: %v, got: %v", c.node,
				c.expected, result)
		}
	}
}

func TestToUnstructured(t *testing.T) {
	t.Skip("Prepare correct test objects")

	cases := []struct {
		node     *kubeapi.Node
		expected *unstructured.Unstructured
	}{
		{
			createTestNode("test-node"),
			&unstructured.Unstructured{},
		},
	}

	for _, c := range cases {
		result, err := ctrl.ToUnstructured(c.node)

		if err != nil {
			t.Errorf("ToUnstructured(node: %v): unexpected error: %v", c.node, err)
		}

		if !reflect.DeepEqual(c.expected, result) {
			t.Errorf("ToUnstructured(node: %v): expected: %v, got: %v", c.node, c.expected, result)
		}
	}
}

func TestToBytes(t *testing.T) {
	t.Skip("Prepare correct test objects")

	cases := []struct {
		unstructured *unstructured.Unstructured
		expected     []byte
	}{
		{
			&unstructured.Unstructured{},
			[]byte{},
		},
	}

	for _, c := range cases {
		result, err := ctrl.ToBytes(c.unstructured)

		if err != nil {
			t.Errorf("ToBytes(unstructured: %v): unexpected error: %v", c.unstructured, err)
		}

		if !reflect.DeepEqual(c.expected, result) {
			t.Errorf("ToBytes(unstructured: %v): expected: %v, got: %v", c.unstructured, c.expected, result)
		}
	}
}
