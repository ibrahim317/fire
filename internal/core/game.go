package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const FPS = 60

type Game struct {
	FontEmoji    rl.Font
	Bg           rl.Texture2D
	GrassBlock   rl.Texture2D
	ScreenWidth  int32
	ScreenHeight int32
	Hero         Character
	Gravity      float32
}

type CharacterState int

const (
	Idle CharacterState = iota
	Running
	Jumping
	Falling
)

type MovementDirection int

const (
	Left MovementDirection = iota
	Right
	Up
	UpRight
	UpLeft
	Down
	DownRight
	DownLeft
)

type AnimationData struct {
	Image        *rl.Image
	Texture      rl.Texture2D
	FrameCount   int32
	CurrentFrame int32
	FrameDelay   int32
	FrameCounter int32
	FrameSize    int32
}

type Character struct {
	States            map[CharacterState]AnimationData
	CurrentState      CharacterState
	MovementDirection MovementDirection
	Position          rl.Vector2
	Velocity          rl.Vector2
	Acceleration      rl.Vector2
}

func (g *Game) Init() {
	g.ScreenWidth = 800
	g.ScreenHeight = 600
	g.Gravity = 0.6
	g.Hero.Velocity = rl.Vector2{X: 0, Y: 0}
	g.Hero.Acceleration = rl.Vector2{X: 0, Y: 0}
	g.Hero.States = make(map[CharacterState]AnimationData)
	g.Hero.CurrentState = Idle
	g.Hero.Position = rl.Vector2{X: 0, Y: 0}
}
