package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/leo9iota/darling-in-the-goland/internal/core"
	"github.com/leo9iota/darling-in-the-goland/internal/entity"
	"github.com/leo9iota/darling-in-the-goland/internal/physics"
	"github.com/leo9iota/darling-in-the-goland/internal/tilemap"
)

const (
	screenWidth  = 640
	screenHeight = 360
	windowWidth  = 1280
	windowHeight = 720
)

// Game implements the ebiten.Game interface.
type Game struct {
	background *core.Background
	tileMap    *tilemap.TileMap
	camera     *core.Camera
	world      *physics.World
	player     *entity.Player
	entities   *entity.Manager
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

	// Create physics world (gravity handled by player/entities, world gravity = 0)
	world := physics.NewWorld(0, 0)

	// Generate static colliders from the solid layer
	colliders := tm.GenerateColliders(world)
	log.Printf("Generated %d colliders from solid layer", len(colliders))

	// Create player
	player, err := entity.NewPlayer(100, 280, world)
	if err != nil {
		return nil, fmt.Errorf("creating player: %w", err)
	}

	// Spawn entities from TMX
	mgr := entity.NewManager()
	mgr.SpawnFromMap(tm.EntitySpawns(), world, player)

	// Create camera
	cam := core.NewCamera(screenWidth, screenHeight)
	cam.SetBounds(tm.MapWidthPx())

	return &Game{
		background: bg,
		tileMap:    tm,
		camera:     cam,
		world:      world,
		player:     player,
		entities:   mgr,
	}, nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	dt := 1.0 / float64(ebiten.TPS())

	// Handle jump input
	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.player.Jump()
	}

	// Update player (movement, gravity, animation)
	g.player.Update(dt)

	// Update entities (coin spin, enemy patrol)
	g.entities.UpdateAll(dt)

	// Step physics (collision detection/resolution, callbacks)
	g.world.Update(dt)

	// Remove collected coins after physics step
	g.entities.RemoveMarked(g.world)

	// Camera follows player
	g.camera.Follow(g.player.Body.Position.X, g.player.Body.Position.Y)
	g.camera.Update(dt)

	// Update background parallax
	g.background.Update(g.camera.X)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 1. Parallax background
	g.background.Draw(screen)

	// 2. Tile layers
	g.tileMap.Draw(screen, g.camera.X, g.camera.Y)

	// 3. Entities (coins, spikes, stones, enemies)
	g.entities.DrawAll(screen, g.camera.X, g.camera.Y)

	// 4. Player (on top of entities)
	g.player.Draw(screen, g.camera.X, g.camera.Y)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
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
