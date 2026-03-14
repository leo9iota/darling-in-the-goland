package main

import (
	"fmt"
	"image"
	_ "image/png" // register PNG decoder
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/leo9iota/darling-in-the-goland/internal/animation"
	"github.com/leo9iota/darling-in-the-goland/internal/core"
	"github.com/leo9iota/darling-in-the-goland/internal/physics"
	"github.com/leo9iota/darling-in-the-goland/internal/tilemap"
)

const (
	screenWidth  = 640
	screenHeight = 360
	windowWidth  = 1280
	windowHeight = 720

	cameraPanSpeed = 150.0 // pixels per second (demo movement)
)

// Game implements the ebiten.Game interface.
type Game struct {
	background *core.Background
	tileMap    *tilemap.TileMap
	camera     *core.Camera
	world      *physics.World
	playerAnim *animation.Animation

	// Demo: simulated player position (moved with arrow keys until real player in Phase 5)
	playerX, playerY float64
}

func NewGame() (*Game, error) {
	// Load parallax background
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

	// Load tilemap
	tm, err := tilemap.New("assets/maps/tmx/map-1.tmx")
	if err != nil {
		return nil, fmt.Errorf("loading tilemap: %w", err)
	}

	// Create physics world (gravity matches Lua: 0, 2000)
	world := physics.NewWorld(0, 2000)

	// Generate static colliders from the solid layer
	colliders := tm.GenerateColliders(world)
	log.Printf("Generated %d colliders from solid layer", len(colliders))

	// Log entity spawns (entities created in Phase 5/6)
	spawns := tm.EntitySpawns()
	for _, sp := range spawns {
		log.Printf("Entity spawn: type=%s pos=(%.0f, %.0f) size=(%.0f, %.0f)",
			sp.Type, sp.X, sp.Y, sp.Width, sp.Height)
	}

	// Load player idle animation
	idleFrames, err := loadFrames("assets/player/idle/zero-two-idle-%d.png", 4)
	if err != nil {
		return nil, fmt.Errorf("loading player idle frames: %w", err)
	}
	playerAnim := animation.NewAnimation(idleFrames, 0.1)

	// Create camera
	cam := core.NewCamera(screenWidth, screenHeight)
	cam.SetBounds(tm.MapWidthPx())

	// Start player near left side of map, on the ground area
	startX := 100.0
	startY := 280.0

	return &Game{
		background: bg,
		tileMap:    tm,
		camera:     cam,
		world:      world,
		playerAnim: playerAnim,
		playerX:    startX,
		playerY:    startY,
	}, nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	dt := 1.0 / float64(ebiten.TPS())

	// Demo: move simulated player position with arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.playerX += cameraPanSpeed * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.playerX -= cameraPanSpeed * dt
	}
	if g.playerX < 0 {
		g.playerX = 0
	}

	// Camera follows the simulated player
	g.camera.Follow(g.playerX, g.playerY)
	g.camera.Update(dt)

	// Update background parallax from camera position
	g.background.Update(g.camera.X)

	g.playerAnim.Update(dt)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 1. Draw parallax background
	g.background.Draw(screen)

	// 2. Draw tile layers
	g.tileMap.Draw(screen, g.camera.X, g.camera.Y)

	// 3. Draw player sprite at its world position, offset by camera
	frame := g.playerAnim.CurrentFrame()
	opts := &ebiten.DrawImageOptions{}
	fw := float64(frame.Bounds().Dx())
	fh := float64(frame.Bounds().Dy())
	opts.GeoM.Translate(g.playerX-fw/2-g.camera.X, g.playerY-fh/2-g.camera.Y)
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
