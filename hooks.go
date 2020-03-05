package hooks

import "sync"

// HookFunc is the type of a hook callback function, receiving a slice of any value as arguments
// and returning one argument and an error.
type HookFunc func(a []interface{}) (interface{}, error)

var lock sync.RWMutex
var cb map[int]HookFunc

// Add adds the given hook for the ID. If a function already exists under that ID, it is overwritten.
func Add(id int, f HookFunc) {
	lock.Lock()
	defer lock.Unlock()
	cb[id] = f
}

// Exec executes the hook for given ID and arguments, if there is any. It does nothing if there is none.
func Exec(id int, args ...interface{}) {
	lock.RLock()
	defer lock.RUnlock()
	if f, ok := cb[id]; ok {
		f(args)
	}
}

// Remove the given hook, if it exists.
func Remove(id int) {
	lock.Lock()
	defer lock.Unlock()
	delete(cb, id)
}

// Active returns true if the given hook is set, false otherwise. Using this first before using Exec
// may be more efficient due to the arguments provided to Exec.
func Active(id int) bool {
	_, ok := cb[id]
	return ok
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	cb = make(map[int]HookFunc)
}
