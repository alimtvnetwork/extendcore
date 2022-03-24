package mappedfunc

import (
	"sync"
)

var (
	wrapErrorWithDetailsInstance = wrapErrorWithDetails{}
	globalMutex                  = sync.Mutex{}
	Converter                    = converter{}
)
