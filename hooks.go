package hooks

import (
	"sync"
)

var lock sync.RWMutex
var hooks map[int]*HookContainer

type HookContainer struct {
	procs     map[int]func(a []interface{})
	keys      []int
	suspended bool
	mutex     sync.RWMutex
	counter   int
}

// newHookContainer creates a new hook function container. Normally, this function does not need to be used.
func newHookContainer() *HookContainer {
	m := make(map[int]func(a []interface{}))
	k := make([]int, 0)
	return &HookContainer{
		procs: m,
		keys:  k,
	}
}

// lock the container and suspend its exec operations.
func (h *HookContainer) lock() {
	h.mutex.Lock()
	h.suspended = true
}

// unlock the container and resume its exec operations.
func (h *HookContainer) unlock() {
	h.mutex.Unlock()
	h.suspended = false
}

// isSuspended returns true if the exec operations are suspended, false otherwise.
func (h *HookContainer) isSuspended() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.suspended
}

// add a new hook to the given container, return its internal id.
func (h *HookContainer) add(f func(a []interface{})) int {
	h.lock()
	defer h.unlock()
	h.counter++
	h.procs[h.counter] = f
	h.keys = append(h.keys, h.counter)
	return h.counter
}

// exec executes the hooks for the given container in LIFO order.
func (h *HookContainer) exec(args []interface{}) {
	h.lock()
	defer h.unlock()
	for i := len(h.keys) - 1; i >= 0; i-- {
		f, ok := h.procs[h.keys[i]]
		if ok {
			f(args)
		}
	}
}

// remove the procedure with id.
func (h *HookContainer) remove(id int) {
	h.lock()
	defer h.unlock()
	delete(h.procs, id)
	for i := range h.keys {
		if h.keys[i] == id {
			h.keys = append(h.keys[:i], h.keys[i+1:]...)
			break
		}
	}
}

// remove all hooks in the container.
func (h *HookContainer) removeAll() {
	h.lock()
	defer h.unlock()
	for k := range h.procs {
		delete(h.procs, k)
	}
	h.keys = make([]int, 0)
}

// Add a function for the given hook. The function is added to the list of functions of the hook.
// Add returns an ID2 that represents the individual function of the hook.
func Add(hook int, f func(a []interface{})) int {
	lock.Lock()
	defer lock.Unlock()
	container, ok := hooks[hook]
	if !ok || container == nil {
		container = newHookContainer()
		hooks[hook] = container
	}
	return container.add(f)
}

// Exec executes all functions for the hook with the given args. It does nothing if there is no function
// for the hook. The functions for a hook are executed in LIFO order. While functions for a hook
// are executed, the hook itself is not called.
func Exec(hook int, args ...interface{}) {
	lock.RLock()
	defer lock.RUnlock()
	container, ok := hooks[hook]
	if !ok || container == nil {
		return
	}
	container.exec(args)
}

// Remove the function with ID for the given hook, if it exists.
func Remove(hook, id int) {
	lock.RLock()
	defer lock.RUnlock()
	container, ok := hooks[hook]
	if !ok || container == nil {
		return
	}
	container.remove(id)
}

// Remove all functions for the hook.
func RemoveAll(hook int) {
	lock.RLock()
	defer lock.RUnlock()
	container, ok := hooks[hook]
	if !ok || container == nil {
		return
	}
	container.removeAll()
	hooks[hook] = nil
}

// Active returns true if the hook is set, false otherwise. Using this first before using Exec
// may be more efficient due to the arguments provided to Exec.
func Active(hook int) bool {
	lock.RLock()
	defer lock.RUnlock()
	container, ok := hooks[hook]
	if !ok || container == nil {
		return false
	}
	return !container.isSuspended()
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	hooks = make(map[int]*HookContainer)
}
