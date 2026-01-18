package core

import (
	"os"
	"path/filepath"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	projectRootOnce sync.Once
	projectRoot     string
)

func (g *Game) LoadAssets() {
	g.FontEmoji = rl.LoadFont(ResourcePath("resources/fonts/dejavu.fnt"))
	g.Bg = rl.LoadTexture(ResourcePath("resources/background/Background.png"))
	g.GrassTile = rl.LoadTexture(ResourcePath("resources/assets/tiles/grass.png"))
	g.HealthHeart = rl.LoadTexture(ResourcePath("resources/assets/heart.png"))

	gifDir := ResourcePath("resources/character/colour2/no_outline/120x80_gifs/")

	// Load mob sprite sheet
	g.MobAnimation = loadSpriteSheetData(ResourcePath("resources/mob/Snail/walk-Sheet.png"), 8, 8, 48*32)

	// Load hero animations
	g.HeroAnimations[int(IdleLegacy)] = loadAnimatedGifData(gifDir+"/__Idle.gif", 8)
	g.HeroAnimations[int(RunningLegacy)] = loadAnimatedGifData(gifDir+"/__Run.gif", 6)
	g.HeroAnimations[int(JumpingLegacy)] = loadAnimatedGifData(gifDir+"/__Jump.gif", 6)
	g.HeroAnimations[int(FallingLegacy)] = loadAnimatedGifData(gifDir+"/__Fall.gif", 6)
}

func (g *Game) UnloadAssets() {
	rl.UnloadFont(g.FontEmoji)
	rl.UnloadTexture(g.Bg)
	rl.UnloadTexture(g.GrassTile)
	rl.UnloadTexture(g.HealthHeart)

	// Unload hero animations
	for _, animData := range g.HeroAnimations {
		rl.UnloadTexture(animData.Texture)
		if animData.Image != nil {
			rl.UnloadImage(animData.Image)
		}
	}

	// Unload mob animation
	rl.UnloadTexture(g.MobAnimation.Texture)
	if g.MobAnimation.Image != nil {
		rl.UnloadImage(g.MobAnimation.Image)
	}
}

// loadAnimatedGifData loads a GIF animation and returns the data.
func loadAnimatedGifData(imagePath string, frameDelay int32) AnimationDataLegacy {
	var frames int32 = 0
	image := rl.LoadImageAnim(imagePath, &frames)
	texture := rl.LoadTextureFromImage(image)
	return AnimationDataLegacy{
		Image:         image,
		Texture:       texture,
		FrameCount:    frames,
		CurrentFrame:  0,
		FrameDelay:    frameDelay,
		FrameCounter:  0,
		FrameSize:     image.Width * image.Height,
		IsSpriteSheet: false,
	}
}

// loadSpriteSheetData loads a sprite sheet and returns the data.
func loadSpriteSheetData(imagePath string, frameCount int32, frameDelay int32, frameSize int32) AnimationDataLegacy {
	image := rl.LoadImage(imagePath)
	texture := rl.LoadTextureFromImage(image)
	return AnimationDataLegacy{
		Image:         image,
		Texture:       texture,
		FrameCount:    frameCount,
		CurrentFrame:  0,
		FrameDelay:    frameDelay,
		FrameCounter:  0,
		FrameSize:     frameSize,
		IsSpriteSheet: true,
	}
}

// ResourcePath returns the absolute path for a resource file.
func ResourcePath(rel string) string {
	root := resolveProjectRoot()
	return filepath.Join(root, rel)
}

func resolveProjectRoot() string {
	projectRootOnce.Do(func() {
		if root, err := findProjectRoot(); err == nil {
			projectRoot = root
		} else {
			projectRoot = "."
		}
	})
	return projectRoot
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, statErr := os.Stat(filepath.Join(dir, "resources")); statErr == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}
