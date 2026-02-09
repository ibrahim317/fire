package core

import (
	"fire/internal/components"
	"fire/internal/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// legacyToAnimationData converts legacy asset data to a component AnimationData.
func legacyToAnimationData(legacy AnimationDataLegacy) *components.AnimationData {
	return &components.AnimationData{
		Image:         legacy.Image,
		Texture:       legacy.Texture,
		FrameCount:    legacy.FrameCount,
		CurrentFrame:  0,
		FrameDelay:    legacy.FrameDelay,
		FrameCounter:  0,
		FrameSize:     legacy.FrameSize,
		IsSpriteSheet: legacy.IsSpriteSheet,
	}
}

// buildHeroAnimations builds the hero animation map from game legacy data.
func buildHeroAnimations(game *Game) map[components.AnimationState]*components.AnimationData {
	animations := make(map[components.AnimationState]*components.AnimationData)
	animations[components.AnimIdle] = legacyToAnimationData(game.GetHeroAnimationData(int(IdleLegacy)))
	animations[components.AnimRunning] = legacyToAnimationData(game.GetHeroAnimationData(int(RunningLegacy)))
	animations[components.AnimJumping] = legacyToAnimationData(game.GetHeroAnimationData(int(JumpingLegacy)))
	animations[components.AnimFalling] = legacyToAnimationData(game.GetHeroAnimationData(int(FallingLegacy)))
	return animations
}

// SpawnPlayer creates the player entity with all required components.
func SpawnPlayer(world *ecs.World, game *Game) *ecs.Entity {
	transformStore := ecs.RegisterStore[*components.TransformComponent](world.Components)
	spriteStore := ecs.RegisterStore[*components.SpriteComponent](world.Components)
	colliderStore := ecs.RegisterStore[*components.ColliderComponent](world.Components)
	inputStore := ecs.RegisterStore[*components.InputComponent](world.Components)
	physicsStore := ecs.RegisterStore[*components.PhysicsComponent](world.Components)
	healthStore := ecs.RegisterStore[*components.HealthComponent](world.Components)

	player := world.CreateEntity("player")

	transformStore.Add(player.ID, &components.TransformComponent{
		Position:    rl.Vector2{X: PlayerStartX, Y: PlayerStartY},
		Velocity:    rl.Vector2{X: 0, Y: 0},
		Acceleration: rl.Vector2{X: 0, Y: 0},
		FacingRight: true,
	})

	animations := buildHeroAnimations(game)
	spriteStore.Add(player.ID, &components.SpriteComponent{
		Animations:  animations,
		CurrentAnim: components.AnimIdle,
		Scale:       game.HeroScaling,
	})

	idleTexture := animations[components.AnimIdle].Texture
	colliderWidth := float32(idleTexture.Width) * game.HeroScaling
	colliderHeight := float32(idleTexture.Height) * game.HeroScaling

	colliderStore.Add(player.ID, &components.ColliderComponent{
		Bounds:    rl.Rectangle{X: 0, Y: 0, Width: colliderWidth, Height: colliderHeight},
		IsTrigger: false,
		Layer:     "player",
	})

	inputStore.Add(player.ID, &components.InputComponent{
		MoveX:       0,
		JumpPressed: false,
		JumpHeld:    false,
	})

	physicsStore.Add(player.ID, &components.PhysicsComponent{
		Gravity:    game.Gravity,
		JumpForce:  PlayerJumpForce,
		MoveSpeed:  PlayerMoveSpeed,
		IsOnGround: false,
	})

	healthStore.Add(player.ID, &components.HealthComponent{
		Current: PlayerHealthMax,
		Max:     PlayerHealthMax,
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

	animData := legacyToAnimationData(game.MobAnimation)
	animations := map[components.AnimationState]*components.AnimationData{
		components.AnimIdle:    animData,
		components.AnimRunning: animData,
	}

	spriteStore.Add(mob.ID, &components.SpriteComponent{
		Animations:  animations,
		CurrentAnim: components.AnimRunning,
		Scale:       1.0,
	})

	colliderStore.Add(mob.ID, &components.ColliderComponent{
		Bounds:    rl.Rectangle{X: 0, Y: 0, Width: MobColliderWidth, Height: MobColliderHeight},
		IsTrigger: false,
		Layer:     "enemy",
	})

	physicsStore.Add(mob.ID, &components.PhysicsComponent{
		Gravity:    game.Gravity,
		JumpForce:  0,
		MoveSpeed:  MobMoveSpeed,
		IsOnGround: false,
	})

	aiStore.Add(mob.ID, &components.AIComponent{
		Behavior:   components.AIPatrol,
		PatrolPath: []rl.Vector2{{X: x, Y: y}, {X: x - MobPatrolDistance, Y: y}},
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
			transform.Position = rl.Vector2{X: PlayerStartX, Y: PlayerStartY}
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
