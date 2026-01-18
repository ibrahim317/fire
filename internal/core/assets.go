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
	g.FontEmoji = rl.LoadFont(resourcePath("resources/fonts/dejavu.fnt"))
	g.Bg = rl.LoadTexture(resourcePath("resources/background/Background.png"))
	g.GrassTile = rl.LoadTexture(resourcePath("resources/assets/tiles/grass.png"))
	g.HealthHeart = rl.LoadTexture(resourcePath("resources/assets/heart.png"))
	gifDir := resourcePath("resources/character/colour2/no_outline/120x80_gifs/")
	loadSpriteSheet(g, resourcePath("resources/mob/Snail/walk-Sheet.png"), 8, 8, 48*32)

	loadAnimatedGif(g, gifDir+"/__Idle.gif", 8, Idle)
	loadAnimatedGif(g, gifDir+"/__Run.gif", 6, Running)
	loadAnimatedGif(g, gifDir+"/__Jump.gif", 6, Jumping)
	loadAnimatedGif(g, gifDir+"/__Fall.gif", 6, Falling)

}

func (g Game) UnloadAssets() {
	rl.UnloadFont(g.FontEmoji)
	rl.UnloadTexture(g.Bg)
	rl.UnloadTexture(g.GrassTile)

	// Unload animated textures and images
	for _, animData := range g.Hero.States {
		rl.UnloadTexture(animData.Texture)
		rl.UnloadImage(animData.Image)
	}
}

// Load animated GIF for Running state
func loadAnimatedGif(g *Game, imagePath string, frameDelay int32, state CharacterState) {
	var frames int32 = 0
	image := rl.LoadImageAnim(imagePath, &frames)
	texture := rl.LoadTextureFromImage(image)
	g.Hero.States[state] = AnimationData{
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

func loadSpriteSheet(g *Game, imagePath string, frameCount int32, frameDelay int32, frameSize int32) {
	image := rl.LoadImage(imagePath)
	texture := rl.LoadTextureFromImage(image)
	g.Mob.AnimationData = AnimationData{
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

func resourcePath(rel string) string {
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
