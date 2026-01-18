package ecs

// EntityID is a unique identifier for an entity.
type EntityID uint32

// Entity represents a game object composed of components.
type Entity struct {
	ID     EntityID
	Active bool
	Tags   []string
}

// HasTag checks if the entity has a specific tag.
func (e *Entity) HasTag(tag string) bool {
	for _, t := range e.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// AddTag adds a tag to the entity if it doesn't already exist.
func (e *Entity) AddTag(tag string) {
	if !e.HasTag(tag) {
		e.Tags = append(e.Tags, tag)
	}
}

// RemoveTag removes a tag from the entity.
func (e *Entity) RemoveTag(tag string) {
	for i, t := range e.Tags {
		if t == tag {
			e.Tags = append(e.Tags[:i], e.Tags[i+1:]...)
			return
		}
	}
}

// EntityManager handles entity creation and ID generation.
type EntityManager struct {
	nextID   EntityID
	entities map[EntityID]*Entity
}

// NewEntityManager creates a new EntityManager.
func NewEntityManager() *EntityManager {
	return &EntityManager{
		nextID:   1,
		entities: make(map[EntityID]*Entity),
	}
}

// CreateEntity creates a new entity with the given tags.
func (em *EntityManager) CreateEntity(tags ...string) *Entity {
	entity := &Entity{
		ID:     em.nextID,
		Active: true,
		Tags:   tags,
	}
	em.entities[em.nextID] = entity
	em.nextID++
	return entity
}

// GetEntity retrieves an entity by ID.
func (em *EntityManager) GetEntity(id EntityID) *Entity {
	return em.entities[id]
}

// RemoveEntity removes an entity by ID.
func (em *EntityManager) RemoveEntity(id EntityID) {
	delete(em.entities, id)
}

// GetAllEntities returns all active entities.
func (em *EntityManager) GetAllEntities() []*Entity {
	result := make([]*Entity, 0, len(em.entities))
	for _, e := range em.entities {
		if e.Active {
			result = append(result, e)
		}
	}
	return result
}

// GetEntitiesWithTag returns all entities with the specified tag.
func (em *EntityManager) GetEntitiesWithTag(tag string) []*Entity {
	result := make([]*Entity, 0)
	for _, e := range em.entities {
		if e.Active && e.HasTag(tag) {
			result = append(result, e)
		}
	}
	return result
}
