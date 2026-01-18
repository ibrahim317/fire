package core

import (
	"fire/internal/components"
	"fire/internal/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// SpawnPlayer creates the player entity with all required components.
func SpawnPlayer(world *ecs.World, game *Game) *ecs.Entity {
	// Register component stores if not already registered
	transformStore := ecs.RegisterStore[*components.TransformComponent](world.Components)
	spriteStore := ecs.RegisterStore[*components.SpriteComponent](world.Components)
	colliderStore := ecs.RegisterStore[*components.ColliderComponent](world.Components)
	inputStore := ecs.RegisterStore[*components.InputComponent](world.Components)
	physicsStore := ecs.RegisterStore[*components.PhysicsComponent](world.Components)
	healthStore := ecs.RegisterStore[*components.HealthComponent](world.Components)

	// Create player entity
	player := world.CreateEntity("player")

	// Add transform component
	transformStore.Add(player.ID, &components.TransformComponent{
		Position:    rl.Vector2{X: 100, Y: 300}, // Start a bit down and right
		Velocity:    rl.Vector2{X: 0, Y: 0},
		Acceleration: rl.Vector2{X: 0, Y: 0},
		FacingRight: true,
	})

	// Convert legacy animation data to component format
	animations := make(map[components.AnimationState]*components.AnimationData)

	idleData := game.GetHeroAnimationData(int(IdleLegacy))
	animations[components.AnimIdle] = &components.AnimationData{
		Image:         idleData.Image,
		Texture:       idleData.Texture,
		FrameCount:    idleData.FrameCount,
		CurrentFrame:  0,
		FrameDelay:    idleData.FrameDelay,
		FrameCounter:  0,
		FrameSize:     idleData.FrameSize,
		IsSpriteSheet: idleData.IsSpriteSheet,
	}

	runData := game.GetHeroAnimationData(int(RunningLegacy))
	animations[components.AnimRunning] = &components.AnimationData{
		Image:         runData.Image,
		Texture:       runData.Texture,
		FrameCount:    runData.FrameCount,
		CurrentFrame:  0,
		FrameDelay:    runData.FrameDelay,
		FrameCounter:  0,
		FrameSize:     runData.FrameSize,
		IsSpriteSheet: runData.IsSpriteSheet,
	}

	jumpData := game.GetHeroAnimationData(int(JumpingLegacy))
	animations[components.AnimJumping] = &components.AnimationData{
		Image:         jumpData.Image,
		Texture:       jumpData.Texture,
		FrameCount:    jumpData.FrameCount,
		CurrentFrame:  0,
		FrameDelay:    jumpData.FrameDelay,
		FrameCounter:  0,
		FrameSize:     jumpData.FrameSize,
		IsSpriteSheet: jumpData.IsSpriteSheet,
	}

	fallData := game.GetHeroAnimationData(int(FallingLegacy))
	animations[components.AnimFalling] = &components.AnimationData{
		Image:         fallData.Image,
		Texture:       fallData.Texture,
		FrameCount:    fallData.FrameCount,
		CurrentFrame:  0,
		FrameDelay:    fallData.FrameDelay,
		FrameCounter:  0,
		FrameSize:     fallData.FrameSize,
		IsSpriteSheet: fallData.IsSpriteSheet,
	}

	// Add sprite component
	spriteStore.Add(player.ID, &components.SpriteComponent{
		Animations:  animations,
		CurrentAnim: components.AnimIdle,
		Scale:       game.HeroScaling,
	})

	// Calculate collider size from idle animation
	idleTexture := animations[components.AnimIdle].Texture
	colliderWidth := float32(idleTexture.Width) * game.HeroScaling
	colliderHeight := float32(idleTexture.Height) * game.HeroScaling

	// Add collider component
	colliderStore.Add(player.ID, &components.ColliderComponent{
		Bounds:    rl.Rectangle{X: 0, Y: 0, Width: colliderWidth, Height: colliderHeight},
		IsTrigger: false,
		Layer:     "player",
	})

	// Add input component
	inputStore.Add(player.ID, &components.InputComponent{
		MoveX:       0,
		JumpPressed: false,
		JumpHeld:    false,
	})

	// Add physics component
	physicsStore.Add(player.ID, &components.PhysicsComponent{
		Gravity:    game.Gravity,
		JumpForce:  12.0, // Increased jump force
		MoveSpeed:  4.0,  // Increased move speed
		IsOnGround: false,
	})

	// Add health component
	healthStore.Add(player.ID, &components.HealthComponent{
		Current: 5,
		Max:     5,
	})

	return player
}

// SpawnMob creates a mob entity.
func SpawnMob(world *ecs.World, game *Game, x, y float32) *ecs.Entity {
	transformStore := ecs.RegisterStore[*components.TransformComponent](world.Components)
	spriteStore := ecs.RegisterStore[*components.SpriteComponent](world.Components)
	colliderStore := ecs.RegisterStore[*components.ColliderComponent](world.Components)
	physicsStore := ecs.RegisterStore[*components.PhysicsComponent](world.Components)
	aiStore := ecs.RegisterStore[*components.AIComponent](world.Components)

	mob := world.CreateEntity("enemy", "mob")

	transformStore.Add(mob.ID, &components.TransformComponent{
		Position:    rl.Vector2{X: x, Y: y},
		Velocity:    rl.Vector2{X: 0, Y: 0},
		Acceleration: rl.Vector2{X: 0, Y: 0},
		FacingRight: false,
	})

	// Mob only has walk animation, map it to running
	animations := make(map[components.AnimationState]*components.AnimationData)
	mobData := game.MobAnimation
	animData := &components.AnimationData{
		Image:         mobData.Image,
		Texture:       mobData.Texture,
		FrameCount:    mobData.FrameCount,
		CurrentFrame:  0,
		FrameDelay:    mobData.FrameDelay,
		FrameCounter:  0,
		FrameSize:     mobData.FrameSize,
		IsSpriteSheet: mobData.IsSpriteSheet,
	}
	animations[components.AnimIdle] = animData
	animations[components.AnimRunning] = animData

	spriteStore.Add(mob.ID, &components.SpriteComponent{
		Animations:  animations,
		CurrentAnim: components.AnimRunning,
		Scale:       1.0,
	})

	// Mob collider (approximate size)
	colliderStore.Add(mob.ID, &components.ColliderComponent{
		Bounds:    rl.Rectangle{X: 0, Y: 0, Width: 48, Height: 32},
		IsTrigger: false,
		Layer:     "enemy",
	})

	physicsStore.Add(mob.ID, &components.PhysicsComponent{
		Gravity:    game.Gravity,
		JumpForce:  0,
		MoveSpeed:  1.0,
		IsOnGround: false,
	})

	aiStore.Add(mob.ID, &components.AIComponent{
		Behavior:   components.AIPatrol,
		PatrolPath: []rl.Vector2{{X: x, Y: y}, {X: x - 100, Y: y}},
		PathIndex:  0,
	})

	return mob
}

// ResetPlayerPosition resets the player's position and state.
func ResetPlayerPosition(world *ecs.World) {
	transformStore, ok := ecs.GetStore[*components.TransformComponent](world.Components)
	if !ok {
		return
	}

	physicsStore, _ := ecs.GetStore[*components.PhysicsComponent](world.Components)

	for _, entity := range world.GetEntitiesWithTag("player") {
		if transform, ok := transformStore.Get(entity.ID); ok {
			transform.Position = rl.Vector2{X: 100, Y: 300}
			transform.Velocity = rl.Vector2{X: 0, Y: 0}
			transform.Acceleration = rl.Vector2{X: 0, Y: 0}
		}

		if physicsStore != nil {
			if physics, ok := physicsStore.Get(entity.ID); ok {
				physics.IsOnGround = false
			}
		}
	}
}
