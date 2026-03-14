package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 360
	windowWidth  = 1280
	windowHeight = 720
)

// Game implements the ebiten.Game interface.
type Game struct{}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Dark blue-gray background
	screen.Fill(color(0x1a, 0x1a, 0x2e, 0xff))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func color(r, g, b, a byte) _color {
	return _color{r, g, b, a}
}

type _color struct{ R, G, B, A byte }

func (c _color) RGBA() (uint32, uint32, uint32, uint32) {
	return uint32(c.R) * 0x101, uint32(c.G) * 0x101, uint32(c.B) * 0x101, uint32(c.A) * 0x101
}

func main() {
	ebiten.SetWindowTitle("Darling in the GoLand")
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetVsyncEnabled(false)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
