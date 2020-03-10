package hooks

import "sync"

var lock sync.RWMutex
var hooks map[int]*HookContainer

type HookContainer struct {
	procs     map[int]func(a []interface{})
	suspended bool
	mutex     sync.RWMutex
	counter   int
}

// newHookContainer creates a new hook function container. Normally, this function does not need to be used.
func newHookContainer() *HookContainer {
	m := make(map[int]func(a []interface{}))
	return &HookContainer{
		procs:     m,
		suspended: false,
	}
}

// add a new hook to the given container, return its internal id.
func (h *HookContainer) add(f func(a []interface{})) int {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.counter++
	h.procs[h.counter] = f
	return h.counter
}

// exec executes the hooks for the given container in unspecified order.
func (h *HookContainer) exec(args ...interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if h.suspended {
		return
	}
	for _, proc := range h.procs {
		proc(args)
	}
}

// suspend the hook container, which means that exec does not do anything if it's called.
func (h *HookContainer) suspend() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.suspended = true
}

// unsuspend the hook container so that future calls to exec will call the procedures stored.
func (h *HookContainer) unsuspend() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.suspended = false
}

// remove the procedure with id.
func (h *HookContainer) remove(id int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.procs, id)
}

// remove all hooks in the container.
func (h *HookContainer) removeAll() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	for k := range h.procs {
		delete(h.procs, k)
	}
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
// for the hook. The functions for a hook are executed in arbitrary order. While functions for a hook
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
	lock.Lock()
	defer lock.Unlock()
	container, ok := hooks[hook]
	if !ok || container == nil {
		return
	}
	container.remove(id)
}

// Remove all functions for the hook.
func RemoveAll(hook int) {
	lock.Lock()
	defer lock.Unlock()
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
	_, ok := hooks[hook]
	return ok
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	hooks = make(map[int]*HookContainer)
}
