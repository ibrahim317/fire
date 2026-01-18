package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const FPS = 60

type Game struct {
	FontEmoji    rl.Font
	Bg           rl.Texture2D
	GrassTile    rl.Texture2D
	HealthHeart  rl.Texture2D
	Mob          Mob
	ScreenWidth  int32
	ScreenHeight int32
	HeroScaling  float32
	Hero         Character
	Map          LevelMap
	Gravity      float32
	// Mode tracking
	Mode     GameMode
	MainMenu *MainMenu
	Designer *Designer
	Settings *SettingsMenu
	// Game settings
	HighlightBorders bool
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
	Image         *rl.Image
	Texture       rl.Texture2D
	FrameCount    int32
	CurrentFrame  int32
	FrameDelay    int32
	FrameCounter  int32
	FrameSize     int32
	IsSpriteSheet bool
}

type Character struct {
	States            map[CharacterState]AnimationData
	CurrentState      CharacterState
	MovementDirection MovementDirection
	Position          rl.Vector2
	Velocity          rl.Vector2
	Acceleration      rl.Vector2
	IsOnGround        bool
}

type Mob struct {
	AnimationData AnimationData
}

func (g *Game) Init() {
	g.ScreenWidth = 800
	g.ScreenHeight = 600
	g.Gravity = 0.6
	g.HeroScaling = 1.7
	g.Hero.Velocity = rl.Vector2{X: 0, Y: 0}
	g.Hero.Acceleration = rl.Vector2{X: 0, Y: 0}
	g.Hero.States = make(map[CharacterState]AnimationData)
	g.Hero.CurrentState = Idle
	g.Hero.Position = rl.Vector2{X: 0, Y: 0}
	g.Mode = ModeMainMenu
}

// InitUI initializes UI components (must be called after window creation)
func (g *Game) InitUI() {
	g.MainMenu = NewMainMenu(g.ScreenWidth, g.ScreenHeight)
	g.Designer = NewDesigner(g)
	g.Settings = NewSettingsMenu(g.ScreenWidth, g.ScreenHeight, &g.HighlightBorders)
}

// ReloadMap reloads the map from disk (for hot-reloading after design changes)
func (g *Game) ReloadMap() {
	g.Map = InitMap(g)
}

// ResetHeroPosition resets the hero to the starting position
func (g *Game) ResetHeroPosition() {
	g.Hero.Position = rl.Vector2{X: 0, Y: 0}
	g.Hero.Velocity = rl.Vector2{X: 0, Y: 0}
	g.Hero.Acceleration = rl.Vector2{X: 0, Y: 0}
	g.Hero.CurrentState = Idle
	g.Hero.IsOnGround = false
}

