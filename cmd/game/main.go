package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/leo9iota/darling-in-the-goland/internal/animation"
	"github.com/leo9iota/darling-in-the-goland/internal/core"
)

const (
	screenWidth  = 640
	screenHeight = 360
	windowWidth  = 1280
	windowHeight = 720

	cameraPanSpeed = 150.0 // pixels per second
)

// Game implements the ebiten.Game interface.
type Game struct {
	background *core.Background
	playerAnim *animation.Animation
	cameraX    float64
}

func NewGame() (*Game, error) {
	// Load parallax background (closest to farthest)
	bg, err := core.NewBackground(
		[]string{
			"assets/background/jungle/jungle-background-1.png",
			"assets/background/jungle/jungle-background-2.png",
			"assets/background/jungle/jungle-background-3.png",
			"assets/background/jungle/jungle-background-4.png",
		},
		[]float64{0.7, 0.5, 0.3, 0.1},
		screenWidth, screenHeight,
	)
	if err != nil {
		return nil, fmt.Errorf("loading background: %w", err)
	}

	// Load player idle animation for demo
	idleFrames, err := loadFrames("assets/player/idle/zero-two-idle-%d.png", 4)
	if err != nil {
		return nil, fmt.Errorf("loading player idle frames: %w", err)
	}
	playerAnim := animation.NewAnimation(idleFrames, 0.1)

	return &Game{
		background: bg,
		playerAnim: playerAnim,
	}, nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	dt := 1.0 / float64(ebiten.TPS())

	// Pan camera with arrow keys for demo
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.cameraX += cameraPanSpeed * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.cameraX -= cameraPanSpeed * dt
	}
	if g.cameraX < 0 {
		g.cameraX = 0
	}

	g.background.Update(g.cameraX)
	g.playerAnim.Update(dt)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw parallax background
	g.background.Draw(screen)

	// Draw animated player sprite at center of screen
	frame := g.playerAnim.CurrentFrame()
	opts := &ebiten.DrawImageOptions{}
	fw := float64(frame.Bounds().Dx())
	fh := float64(frame.Bounds().Dy())
	opts.GeoM.Translate(float64(screenWidth)/2-fw/2, float64(screenHeight)/2-fh/2)
	screen.DrawImage(frame, opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// loadFrames loads numbered PNG frames using a format pattern (e.g. "dir/sprite-%d.png").
func loadFrames(pattern string, count int) ([]*ebiten.Image, error) {
	frames := make([]*ebiten.Image, 0, count)
	for i := 1; i <= count; i++ {
		path := fmt.Sprintf(pattern, i)
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("opening %s: %w", path, err)
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("decoding %s: %w", path, err)
		}
		frames = append(frames, ebiten.NewImageFromImage(img))
	}
	return frames, nil
}

func main() {
	ebiten.SetWindowTitle("Darling in the GoLand")
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetVsyncEnabled(false)

	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
