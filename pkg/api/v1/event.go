package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

// Event represents a single event to a watched resource.
// This is copy of watch.Event object but it contains annotations for json serialization needed
// for sending request back to kubelet.
type Event struct {
	Type watch.EventType `json:"type"`

	// Object is:
	//  * If Type is Added or Modified: the new state of the object.
	//  * If Type is Deleted: the state of the object immediately before deletion.
	//  * If Type is Error: *api.Status is recommended; other types may make sense
	//    depending on context.
	Object runtime.Object `json:"object"`
}
