package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// BackgroundScale is the scale factor used when drawing fullscreen backgrounds.
const BackgroundScale = 2.7

// DrawScreenBackground draws a fullscreen background texture and an optional dark overlay.
func DrawScreenBackground(bg rl.Texture2D, overlayAlpha uint8) {
	rl.DrawTextureEx(bg, rl.Vector2{X: 0, Y: 0}, 0, BackgroundScale, rl.White)
	w := int32(rl.GetScreenWidth())
	h := int32(rl.GetScreenHeight())
	rl.DrawRectangle(0, 0, w, h, rl.Color{R: 0, G: 0, B: 0, A: overlayAlpha})
}

// DrawScreenTitle draws a centered title at the given y position.
func DrawScreenTitle(title string, fontSize int32, y int32, color rl.Color) {
	w := int32(rl.GetScreenWidth())
	titleWidth := rl.MeasureText(title, fontSize)
	x := w/2 - titleWidth/2
	rl.DrawText(title, x, y, fontSize, color)
}

// DrawScreenTitleWithShadow draws a centered title with a drop shadow.
func DrawScreenTitleWithShadow(title string, fontSize int32, y int32, color rl.Color) {
	w := int32(rl.GetScreenWidth())
	titleWidth := rl.MeasureText(title, fontSize)
	x := w/2 - titleWidth/2
	shadow := rl.Color{R: 0, G: 0, B: 0, A: 180}
	rl.DrawText(title, x+3, y+3, fontSize, shadow)
	rl.DrawText(title, x, y, fontSize, color)
}

// DrawScreenInstructions draws centered instruction text at the given y (e.g. bottom of screen).
func DrawScreenInstructions(text string, y int32, fontSize int32, color rl.Color) {
	w := int32(rl.GetScreenWidth())
	textWidth := rl.MeasureText(text, fontSize)
	rl.DrawText(text, w/2-textWidth/2, y, fontSize, color)
}
