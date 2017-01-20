package common

import (
	"sync"

	"github.com/pborman/uuid"
	"k8s.io/apimachinery/pkg/types"
)

var uuidLock sync.Mutex
var lastUUID uuid.UUID

func NewUUID() types.UID {
	uuidLock.Lock()
	defer uuidLock.Unlock()
	result := uuid.NewUUID()

	// The UUID package is naive and can generate identical UUIDs if the time interval is
	// quick enough. The UUID uses 100 ns increments so it's short enough to actively wait for
	// a new value.
	for uuid.Equal(lastUUID, result) == true {
		result = uuid.NewUUID()
	}

	lastUUID = result
	return types.UID(result.String())
}
