package systems

import (
	"fire/internal/components"
	"fire/internal/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// CollisionSystem detects and resolves collisions between entities.
type CollisionSystem struct{}

// NewCollisionSystem creates a new CollisionSystem.
func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{}
}

// Update checks and resolves collisions for all collidable entities.
func (s *CollisionSystem) Update(world *ecs.World, dt float32) {
	transformStore, ok1 := ecs.GetStore[*components.TransformComponent](world.Components)
	colliderStore, ok2 := ecs.GetStore[*components.ColliderComponent](world.Components)
	physicsStore, _ := ecs.GetStore[*components.PhysicsComponent](world.Components)
	tileStore, _ := ecs.GetStore[*components.TileComponent](world.Components)

	if !ok1 || !ok2 {
		return
	}

	// Get all tile entities for collision checking
	tileEntities := make([]*ecs.Entity, 0)
	if tileStore != nil {
		for _, id := range tileStore.All() {
			entity := world.GetEntity(id)
			if entity != nil && entity.Active {
				tileEntities = append(tileEntities, entity)
			}
		}
	}

	// Check collisions for each non-tile entity
	for _, id := range colliderStore.All() {
		// Skip tiles (they don't move)
		if tileStore != nil && tileStore.Has(id) {
			continue
		}

		entity := world.GetEntity(id)
		if entity == nil || !entity.Active {
			continue
		}

		if !transformStore.Has(id) {
			continue
		}

		transform, _ := transformStore.Get(id)
		collider, _ := colliderStore.Get(id)

		// Reset ground state
		hasPhysics := physicsStore != nil && physicsStore.Has(id)
		if hasPhysics {
			physics, _ := physicsStore.Get(id)
			physics.IsOnGround = false
		}

		// First pass: Resolve vertical collisions
		for _, tileEntity := range tileEntities {
			tileTransform, ok := transformStore.Get(tileEntity.ID)
			if !ok {
				continue
			}
			tileCollider, ok := colliderStore.Get(tileEntity.ID)
			if !ok {
				continue
			}

			entityBounds := collider.GetWorldBounds(transform.Position)
			tileBounds := tileCollider.GetWorldBounds(tileTransform.Position)

			collisionDir := checkCollisionDirection(entityBounds, tileBounds)
			if collisionDir.Y != 0 {
				collisionRec := rl.GetCollisionRec(entityBounds, tileBounds)
				correctionY := collisionDir.Y * collisionRec.Height
				transform.Position.Y += correctionY

				if collisionDir.Y*transform.Velocity.Y < 0 {
					transform.Velocity.Y = 0
				}

				if collisionDir.Y == -1 && hasPhysics {
					physics, _ := physicsStore.Get(id)
					physics.IsOnGround = true
				}
			}
		}

		// Second pass: Resolve horizontal collisions
		for _, tileEntity := range tileEntities {
			tileTransform, ok := transformStore.Get(tileEntity.ID)
			if !ok {
				continue
			}
			tileCollider, ok := colliderStore.Get(tileEntity.ID)
			if !ok {
				continue
			}

			entityBounds := collider.GetWorldBounds(transform.Position)
			tileBounds := tileCollider.GetWorldBounds(tileTransform.Position)

			collisionDir := checkCollisionDirection(entityBounds, tileBounds)
			if collisionDir.X != 0 {
				collisionRec := rl.GetCollisionRec(entityBounds, tileBounds)
				correctionX := collisionDir.X * collisionRec.Width
				transform.Position.X += correctionX

				if collisionDir.X*transform.Velocity.X < 0 {
					transform.Velocity.X = 0
				}
			}
		}
	}
}

// checkCollisionDirection returns the direction of collision between two rectangles.
func checkCollisionDirection(r1, r2 rl.Rectangle) rl.Vector2 {
	r1Left, r1Top, r1Right, r1Bottom := getEdges(r1)
	r2Left, r2Top, r2Right, r2Bottom := getEdges(r2)

	// No collision
	if r1Left > r2Right || r1Right < r2Left || r1Bottom < r2Top || r1Top > r2Bottom {
		return rl.Vector2{X: 0, Y: 0}
	}

	// Calculate overlap on each side
	overlapLeft := r1Right - r2Left
	overlapRight := r2Right - r1Left
	overlapTop := r1Bottom - r2Top
	overlapBottom := r2Bottom - r1Top

	// Find the minimum overlap direction
	minOverlap := overlapLeft
	unitVec := rl.Vector2{X: -1, Y: 0} // left

	if overlapRight < minOverlap {
		minOverlap = overlapRight
		unitVec = rl.Vector2{X: 1, Y: 0} // right
	}
	// Prefer vertical collision resolution to avoid "ghost walls"
	if overlapTop <= minOverlap {
		minOverlap = overlapTop
		unitVec = rl.Vector2{X: 0, Y: -1} // top (standing on ground)
	}
	if overlapBottom < minOverlap {
		unitVec = rl.Vector2{X: 0, Y: 1} // bottom (hitting ceiling)
	}

	return unitVec
}

// getEdges returns the edges of a rectangle: left, top, right, bottom.
func getEdges(r rl.Rectangle) (float32, float32, float32, float32) {
	return r.X, r.Y, r.X + r.Width, r.Y + r.Height
}
