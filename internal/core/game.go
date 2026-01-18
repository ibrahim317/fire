package core

import (
	"fire/internal/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const FPS = 60

// Game holds game configuration, assets, UI components, and the ECS world.
type Game struct {
	// Screen settings
	ScreenWidth  int32
	ScreenHeight int32

	// Assets
	FontEmoji   rl.Font
	Bg          rl.Texture2D
	GrassTile   rl.Texture2D
	HealthHeart rl.Texture2D

	// Hero animation data (for spawning player entity)
	HeroAnimations map[int]AnimationDataLegacy
	HeroScaling    float32

	// Mob animation data (for spawning mob entities)
	MobAnimation AnimationDataLegacy

	// Mode tracking
	Mode     GameMode
	MainMenu *MainMenu
	Designer *Designer
	Settings *SettingsMenu

	// Game settings
	HighlightBorders bool
	Gravity          float32

	// ECS World (used in game mode)
	World *ecs.World
}

// AnimationDataLegacy holds animation data for asset loading.
// This is used during asset loading and converted to components when spawning entities.
type AnimationDataLegacy struct {
	Image         *rl.Image
	Texture       rl.Texture2D
	FrameCount    int32
	CurrentFrame  int32
	FrameDelay    int32
	FrameCounter  int32
	FrameSize     int32
	IsSpriteSheet bool
}

// CharacterStateLegacy is used for asset loading keys.
type CharacterStateLegacy int

const (
	IdleLegacy CharacterStateLegacy = iota
	RunningLegacy
	JumpingLegacy
	FallingLegacy
)

func (g *Game) Init() {
	g.ScreenWidth = 800
	g.ScreenHeight = 600
	g.Gravity = 0.6
	g.HeroScaling = 1.7
	g.HeroAnimations = make(map[int]AnimationDataLegacy)
	g.Mode = ModeMainMenu
}

// InitUI initializes UI components (must be called after window creation).
func (g *Game) InitUI() {
	g.MainMenu = NewMainMenu(g.ScreenWidth, g.ScreenHeight)
	g.Designer = NewDesigner(g)
	g.Settings = NewSettingsMenu(g.ScreenWidth, g.ScreenHeight, &g.HighlightBorders)
}

// InitWorld creates a new ECS world for gameplay.
func (g *Game) InitWorld() {
	g.World = ecs.NewWorld()
}

// GetHeroAnimationData returns animation data for a character state.
func (g *Game) GetHeroAnimationData(state int) AnimationDataLegacy {
	return g.HeroAnimations[state]
}
