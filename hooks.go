package hooks

import "sync"

// HookFunc is the type of a hook callback function, receiving a slice of any value as arguments
// and returning one argument as interface{}.
type HookFunc func(a []interface{}) interface{}

var lock sync.RWMutex
var cb map[int]map[int]HookFunc
var counter int

// Add a function for the given hook. The function is added to the list of functions of the hook.
// Add returns an ID2 that represents the individual function of the hook.
func Add(hook int, f HookFunc) int {
	lock.Lock()
	defer lock.Unlock()
	hooks, ok := cb[hook]
	if !ok || hooks == nil {
		hooks = make(map[int]HookFunc)
	}
	counter++
	hooks[counter] = f
	cb[hook] = hooks
	return counter
}

// Exec executes all functions for the hook with the given args. It does nothing if there is no function
// for the hook. The functions for a hook are executed in arbitrary order.
func Exec(hook int, args ...interface{}) {
	lock.RLock()
	defer lock.RUnlock()
	hooks, ok := cb[hook]
	if !ok || hooks == nil {
		return
	}
	for _, f := range hooks {
		f(args)
	}
}

// Remove the function with ID for the given hook, if it exists.
func Remove(hook, id int) {
	lock.Lock()
	defer lock.Unlock()
	hooks, ok := cb[hook]
	if !ok || hooks == nil {
		return
	}
	delete(hooks, id)
}

// Remove all functions for the hook.
func RemoveAll(hook int) {
	lock.Lock()
	defer lock.Unlock()
	delete(cb, hook)
}

// Active returns true if the hook is set, false otherwise. Using this first before using Exec
// may be more efficient due to the arguments provided to Exec.
func Active(hook int) bool {
	lock.RLock()
	defer lock.RUnlock()
	_, ok := cb[hook]
	return ok
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	cb = make(map[int]map[int]HookFunc)
}
