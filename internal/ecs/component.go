package ecs

import "reflect"

// Component is a marker interface for all component types.
type Component interface{}

// ComponentStore holds all components of a specific type, keyed by EntityID.
type ComponentStore[T Component] struct {
	components map[EntityID]T
}

// NewComponentStore creates a new component store.
func NewComponentStore[T Component]() *ComponentStore[T] {
	return &ComponentStore[T]{
		components: make(map[EntityID]T),
	}
}

// Add adds a component for an entity.
func (cs *ComponentStore[T]) Add(id EntityID, component T) {
	cs.components[id] = component
}

// Get retrieves a component for an entity.
func (cs *ComponentStore[T]) Get(id EntityID) (T, bool) {
	c, ok := cs.components[id]
	return c, ok
}

// Remove removes a component for an entity.
func (cs *ComponentStore[T]) Remove(id EntityID) {
	delete(cs.components, id)
}

// Has checks if an entity has this component.
func (cs *ComponentStore[T]) Has(id EntityID) bool {
	_, ok := cs.components[id]
	return ok
}

// All returns all entity IDs that have this component.
func (cs *ComponentStore[T]) All() []EntityID {
	ids := make([]EntityID, 0, len(cs.components))
	for id := range cs.components {
		ids = append(ids, id)
	}
	return ids
}

// ComponentRegistry provides a type-safe way to store multiple component types.
type ComponentRegistry struct {
	stores map[reflect.Type]interface{}
}

// NewComponentRegistry creates a new component registry.
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		stores: make(map[reflect.Type]interface{}),
	}
}

// RegisterStore registers a component store for a specific type.
func RegisterStore[T Component](cr *ComponentRegistry) *ComponentStore[T] {
	var zero T
	t := reflect.TypeOf(zero)
	if store, ok := cr.stores[t]; ok {
		return store.(*ComponentStore[T])
	}
	store := NewComponentStore[T]()
	cr.stores[t] = store
	return store
}

// GetStore retrieves a component store for a specific type.
func GetStore[T Component](cr *ComponentRegistry) (*ComponentStore[T], bool) {
	var zero T
	t := reflect.TypeOf(zero)
	store, ok := cr.stores[t]
	if !ok {
		return nil, false
	}
	return store.(*ComponentStore[T]), true
}
