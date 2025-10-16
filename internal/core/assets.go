package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) LoadAssets() {
	g.FontEmoji = rl.LoadFont("resources/fonts/dejavu.fnt")
	g.Bg = rl.LoadTexture("resources/background/Background.png")
	g.GrassBlock = rl.LoadTexture("resources/assets/blocks/grass.png")

	// Load animated GIF for Idle state
	var idleFrames int32 = 0
	idleImage := rl.LoadImageAnim("resources/character/colour2/no_outline/120x80_gifs/__Idle.gif", &idleFrames)
	idleTexture := rl.LoadTextureFromImage(idleImage)
	g.Hero.States[Idle] = AnimationData{
		Image:        idleImage,
		Texture:      idleTexture,
		FrameCount:   idleFrames,
		CurrentFrame: 0,
		FrameDelay:   8, // Adjust this value to control animation speed
		FrameCounter: 0,
		FrameSize:    idleImage.Width * idleImage.Height,
	}

	// Load animated GIF for Running state
	var runningFrames int32 = 0
	runningImage := rl.LoadImageAnim("resources/character/colour2/no_outline/120x80_gifs/__Run.gif", &runningFrames)
	runningTexture := rl.LoadTextureFromImage(runningImage)
	g.Hero.States[Running] = AnimationData{
		Image:        runningImage,
		Texture:      runningTexture,
		FrameCount:   runningFrames,
		CurrentFrame: 0,
		FrameDelay:   6, // Running animation can be faster
		FrameCounter: 0,
		FrameSize:    runningImage.Width * runningImage.Height,
	}
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
