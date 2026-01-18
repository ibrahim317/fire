package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// TransformComponent holds position and physics state.
type TransformComponent struct {
	Position    rl.Vector2
	Velocity    rl.Vector2
	Acceleration rl.Vector2
	FacingRight bool
}

// AnimationState represents the current animation state.
type AnimationState int

const (
	AnimIdle AnimationState = iota
	AnimRunning
	AnimJumping
	AnimFalling
)

// AnimationData holds data for a single animation.
type AnimationData struct {
	Image         *rl.Image
	Texture       rl.Texture2D
	FrameCount    int32
	CurrentFrame  int32
	FrameDelay    int32
	FrameCounter  int32
	FrameSize     int32
	IsSpriteSheet bool
}

// SpriteComponent holds visual representation data.
type SpriteComponent struct {
	Animations   map[AnimationState]*AnimationData
	CurrentAnim  AnimationState
	Scale        float32
}

// GetCurrentAnimation returns the current animation data.
func (s *SpriteComponent) GetCurrentAnimation() *AnimationData {
	return s.Animations[s.CurrentAnim]
}

// ColliderComponent holds collision detection data.
type ColliderComponent struct {
	// Bounds is relative to the entity's position
	Bounds    rl.Rectangle
	IsTrigger bool   // If true, detects overlap but doesn't block
	Layer     string // "player", "enemy", "ground", "collectible"
}

// GetWorldBounds returns the collider bounds in world coordinates.
func (c *ColliderComponent) GetWorldBounds(position rl.Vector2) rl.Rectangle {
	return rl.Rectangle{
		X:      position.X + c.Bounds.X,
		Y:      position.Y + c.Bounds.Y,
		Width:  c.Bounds.Width,
		Height: c.Bounds.Height,
	}
}

// InputComponent holds player input state.
type InputComponent struct {
	MoveX       float32 // -1.0 to 1.0
	JumpPressed bool    // True on the frame jump was pressed
	JumpHeld    bool    // True while jump key is held
}

// PhysicsComponent holds physics simulation data.
type PhysicsComponent struct {
	Gravity    float32
	JumpForce  float32
	MoveSpeed  float32
	IsOnGround bool
}

// HealthComponent holds health/damage data.
type HealthComponent struct {
	Current int
	Max     int
}

// IsDead returns true if current health is zero or less.
func (h *HealthComponent) IsDead() bool {
	return h.Current <= 0
}

// TakeDamage reduces health by the specified amount.
func (h *HealthComponent) TakeDamage(amount int) {
	h.Current -= amount
	if h.Current < 0 {
		h.Current = 0
	}
}

// Heal increases health by the specified amount, up to max.
func (h *HealthComponent) Heal(amount int) {
	h.Current += amount
	if h.Current > h.Max {
		h.Current = h.Max
	}
}

// TileType represents the type of tile.
type TileType int32

const (
	TileGrass TileType = iota
	TileStone
	TileWater
	TileTree
	TileRock
)

// TileComponent marks an entity as a static tile.
type TileComponent struct {
	TileType TileType
}

// AIBehavior represents the type of AI behavior.
type AIBehavior int

const (
	AIPatrol AIBehavior = iota
	AIChase
	AIIdle
)

// AIComponent holds AI behavior data.
type AIComponent struct {
	Behavior   AIBehavior
	PatrolPath []rl.Vector2
	PathIndex  int
	TargetID   uint32 // EntityID of target (0 if none)
}
