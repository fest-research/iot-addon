package controller

import (
	"k8s.io/apimachinery/pkg/watch"
)

type Controller interface {
	PodController() IPodController
	NodeController() INodeController
}

// Every resource controller that watches for changes has to implement this interface
type WatchEventController interface {
	TransformWatchEvent(watch.Event) watch.Event
}
