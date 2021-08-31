package mutex

import (
	"runtime"
	"sync/atomic"

	"github.com/yehan2002/syncx/internal"
)

//SpinLock a simple spin lock
type SpinLock struct {
	v uint32
	_ internal.NoCopy
}

//Lock locks the spin lock
func (sl *SpinLock) Lock() {
	for !atomic.CompareAndSwapUint32(&sl.v, 0, 1) {
		runtime.Gosched() //without this it locks up on GOMAXPROCS > 1
	}
}

// Unlock unlocks the spin lock.
// Calling this when the lock is unlocked is a noop.
func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl.v, 0)
}
