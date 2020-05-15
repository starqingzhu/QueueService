package connManager

import (
	"sync"
)

var ConnManager = new(sync.Map)
