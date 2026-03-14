package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/leo9iota/darling-in-the-goland/internal/core"
	"github.com/leo9iota/darling-in-the-goland/internal/entity"
	"github.com/leo9iota/darling-in-the-goland/internal/gui"
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
	hud        *gui.HUD
	menu       *gui.Menu
	debug      *gui.Debug
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

	// World gravity matches playerGravity (1500). Player/enemy have GravityScale=0
	// and handle gravity manually. Stones use GravityScale=1 to fall naturally.
	world := physics.NewWorld(0, 1500)

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

	// Create GUI layers
	hud, err := gui.NewHUD(screenWidth, screenHeight)
	if err != nil {
		return nil, fmt.Errorf("creating HUD: %w", err)
	}

	g := &Game{
		background: bg,
		tileMap:    tm,
		camera:     cam,
		world:      world,
		player:     player,
		entities:   mgr,
		hud:        hud,
	}

	// Create menu (needs reference to toggle function)
	menu, err := gui.NewMenu(screenWidth, screenHeight, func() {
		g.menu.Toggle()
	})
	if err != nil {
		return nil, fmt.Errorf("creating menu: %w", err)
	}
	g.menu = menu

	// Create debug overlay
	debug, err := gui.NewDebug()
	if err != nil {
		return nil, fmt.Errorf("creating debug: %w", err)
	}
	g.debug = debug

	return g, nil
}

func (g *Game) Update() error {
	// F3 toggles debug overlay
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		g.debug.Toggle()
	}

	// Escape toggles pause menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.menu.Toggle()
	}

	// Handle mouse clicks for menu
	if g.menu.Active {
		g.menu.Update()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			if err := g.menu.MousePressed(mx, my); err != nil {
				return err
			}
		}
		return nil // Block game updates while menu is open
	}

	dt := 1.0 / float64(ebiten.TPS())

	// Handle jump input
	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.player.Jump()
	}

	// Update player
	g.player.Update(dt)

	// Update entities
	g.entities.UpdateAll(dt)

	// Step physics
	g.world.Update(dt)

	// Remove collected coins
	g.entities.RemoveMarked(g.world)

	// Camera follows player
	g.camera.Follow(g.player.Body.Position.X, g.player.Body.Position.Y)
	g.camera.Update(dt)

	// Update background parallax
	g.background.Update(g.camera.X)

	// Update debug overlay
	g.debug.Update(dt)
	coins, enemies, spikes, stones := g.entities.GetCounts()
	g.debug.UpdateEntityCounts(coins, enemies, spikes, stones)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 1. Parallax background
	g.background.Draw(screen)

	// 2. Tile layers
	g.tileMap.Draw(screen, g.camera.X, g.camera.Y)

	// 3. Entities
	g.entities.DrawAll(screen, g.camera.X, g.camera.Y)

	// 4. Player
	g.player.Draw(screen, g.camera.X, g.camera.Y)

	// 5. HUD (always visible, on top of game)
	g.hud.Draw(screen, g.player.Health(), 3, g.player.CoinCount())

	// 6. Debug overlay (if active)
	g.debug.Draw(screen)

	// 7. Pause menu (on top of everything)
	g.menu.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowTitle("Darling in the GoLand")
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetScreenFilterEnabled(false) // nearest-neighbor for crisp pixel art
	ebiten.SetVsyncEnabled(false)

	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
