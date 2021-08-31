package mutex

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

//LockTimeout the time out for the lock
var LockTimeout = time.Second * 5

//DebugRWMutex a RWMutex that prints locks and unlocks
type DebugRWMutex struct {
	sync.RWMutex
	rLocks  uint64
	wLocker string

	history  map[string]int
	histlock sync.Mutex
}

//Lock lock for writing
func (r *DebugRWMutex) Lock() {
	r.awaitLock(func(c chan struct{}) {
		r.RWMutex.Lock()
		close(c)
	})
	r.wLocker = r.caller("Lock")
}

//Unlock unlock
func (r *DebugRWMutex) Unlock() {
	r.caller("Unlock")
	r.wLocker = ""
	r.RWMutex.Unlock()
}

//RLock lock for reading
func (r *DebugRWMutex) RLock() {
	r.awaitLock(func(c chan struct{}) {
		r.RWMutex.RLock()
		close(c)
	})

	r.caller("RLock")
	atomic.AddUint64(&r.rLocks, 1)
}

//RUnlock unlock
func (r *DebugRWMutex) RUnlock() {
	r.caller("RUnlock")
	r.RWMutex.RUnlock()
	atomic.AddUint64(&r.rLocks, ^uint64(0))
}

func (r *DebugRWMutex) caller(act string) string {
	var holder string
	if _, file, line, ok := runtime.Caller(2); ok {
		holder = fmt.Sprintf("%s: %s:%d", act, file, line)
	}
	//fmt.Println(holder)
	r.histlock.Lock()
	if r.history == nil {
		r.history = map[string]int{}
	}
	if n, ok := r.history[holder]; ok {
		r.history[holder] = n + 1
	} else {
		r.history[holder] = 1
	}
	r.histlock.Unlock()
	return holder
}

func (r *DebugRWMutex) awaitLock(f func(chan struct{})) {
	timeout := time.After(LockTimeout)
	c := make(chan struct{}, 0)
	go f(c)
	select {
	case <-c:
	case <-timeout:
		r.panicTimeout()
	}
}

func (r *DebugRWMutex) panicTimeout() {
	var msg strings.Builder
	msg.WriteString("lock timed out ")
	fmt.Fprintf(&msg, "after %s", LockTimeout)
	if r.wLocker != "" {
		msg.WriteString("\nWrite lock held by:\n")
		msg.WriteString(r.wLocker)
	}
	if r := atomic.LoadUint64(&r.rLocks); r != 0 {
		fmt.Fprintf(&msg, "\nRead lock held %d times", r)
	}
	r.histlock.Lock()
	if len(r.history) != 0 {
		msg.WriteString("Lock history: \n")
		for k, v := range r.history {
			fmt.Fprintf(&msg, "%s: %d\n", k, v)
		}
	}
	r.histlock.Unlock()

	panic(msg.String())
}

//NewDebugRWMutex create a new debug rw mutex
func NewDebugRWMutex() *DebugRWMutex { return new(DebugRWMutex) }

//NewDebugMutex create a new debug mutex
func NewDebugMutex() sync.Locker { return new(DebugRWMutex) }
