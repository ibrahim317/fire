package ecs

// System is an interface for game systems that process entities.
type System interface {
	Update(world *World, dt float32)
}

// World is the main container for all entities, components, and systems.
type World struct {
	Entities   *EntityManager
	Components *ComponentRegistry
	Systems    []System
	Events     *EventBus
}

// NewWorld creates a new World instance.
func NewWorld() *World {
	return &World{
		Entities:   NewEntityManager(),
		Components: NewComponentRegistry(),
		Systems:    make([]System, 0),
		Events:     NewEventBus(),
	}
}

// AddSystem adds a system to the world.
func (w *World) AddSystem(s System) {
	w.Systems = append(w.Systems, s)
}

// Update runs all systems in order.
func (w *World) Update(dt float32) {
	for _, s := range w.Systems {
		s.Update(w, dt)
	}
}

// CreateEntity creates a new entity with the given tags.
func (w *World) CreateEntity(tags ...string) *Entity {
	return w.Entities.CreateEntity(tags...)
}

// RemoveEntity removes an entity and all its components.
func (w *World) RemoveEntity(id EntityID) {
	w.Entities.RemoveEntity(id)
}

// GetEntity retrieves an entity by ID.
func (w *World) GetEntity(id EntityID) *Entity {
	return w.Entities.GetEntity(id)
}

// GetEntitiesWithTag returns all entities with the specified tag.
func (w *World) GetEntitiesWithTag(tag string) []*Entity {
	return w.Entities.GetEntitiesWithTag(tag)
}

// GetAllEntities returns all active entities.
func (w *World) GetAllEntities() []*Entity {
	return w.Entities.GetAllEntities()
}
