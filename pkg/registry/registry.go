package registry

import (
	"log"
	"reflect"
	"sync"
)

// Registry is a registry of types, with an
// instance assigned to it. This is so that
// we can call types by their name.
type Registry struct {
	types map[string]reflect.Type
	mu    sync.Mutex
}

// NewRegistry - returns a new registry instance
func NewRegistry() *Registry {
	return &Registry{
		types: make(map[string]reflect.Type),
		mu:    sync.Mutex{},
	}
}

// Register a type with a string value
func (r *Registry) Register(name string, element interface{}) {
	log.Println("TeSTST")
	r.mu.Lock()
	defer r.mu.Unlock()
	r.types[name] = reflect.TypeOf(element)
	log.Println(r.types)
}

// Deregister a type
func (r *Registry) Deregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.types, name)
}

// MakeInstance of a type by its name
func (r *Registry) MakeInstance(name string) interface{} {
	v := reflect.New(r.types[name])
	return v.Interface()
}
