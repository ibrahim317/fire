package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) LoadAssets() {
	g.FontEmoji = rl.LoadFont("resources/fonts/dejavu.fnt")
	g.Bg = rl.LoadTexture("resources/background/Background.png")
	g.GrassBlock = rl.LoadTexture("resources/assets/blocks/grass.png")

	gifDir := "resources/character/colour2/no_outline/120x80_gifs/"

	loadAnimatedGif(g, gifDir+"__Idle.gif", 8, Idle)
	loadAnimatedGif(g, gifDir+"__Run.gif", 6, Running)
	loadAnimatedGif(g, gifDir+"__Jump.gif", 6, Jumping)
	loadAnimatedGif(g, gifDir+"__Fall.gif", 6, Falling)

}

func (g Game) UnloadAssets() {
	rl.UnloadFont(g.FontEmoji)
	rl.UnloadTexture(g.Bg)
	rl.UnloadTexture(g.GrassBlock)

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
		Image:        image,
		Texture:      texture,
		FrameCount:   frames,
		CurrentFrame: 0,
		FrameDelay:   frameDelay,
		FrameCounter: 0,
		FrameSize:    image.Width * image.Height,
	}
}
