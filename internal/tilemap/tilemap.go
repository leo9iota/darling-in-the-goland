package tilemap

import (
	"fmt"

	ebiten "github.com/hajimehoshi/ebiten/v2"

	"github.com/leo9iota/darling-in-the-goland/internal/physics"
)

// TileMap is the high-level map struct used for rendering and physics integration.
type TileMap struct {
	parsed  *ParsedMap
	tileset *Tileset

	// renderOrder defines which layers to draw and in what order.
	renderOrder []string
}

// New loads a TMX map and its tileset, returning a ready-to-use TileMap.
func New(tmxPath string) (*TileMap, error) {
	parsed, err := LoadMap(tmxPath)
	if err != nil {
		return nil, fmt.Errorf("loading map: %w", err)
	}

	ts, err := LoadTileset(
		parsed.TilesetImagePath,
		parsed.TilesetFirstGID,
		parsed.TileWidth,
		parsed.TileHeight,
		parsed.TilesetColumns,
	)
	if err != nil {
		return nil, fmt.Errorf("loading tileset: %w", err)
	}

	// Render ground first, then grass on top. Skip solid/entity (invisible).
	order := []string{}
	for _, name := range []string{"ground", "grass"} {
		if _, ok := parsed.Layers[name]; ok {
			order = append(order, name)
		}
	}

	return &TileMap{
		parsed:      parsed,
		tileset:     ts,
		renderOrder: order,
	}, nil
}

// Draw renders the visible tile layers offset by the camera position.
// Only tiles within the viewport are drawn for performance.
func (tm *TileMap) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	tw, th := tm.parsed.TileWidth, tm.parsed.TileHeight

	// Compute visible tile range
	startCol := int(cameraX) / tw
	startRow := int(cameraY) / th
	endCol := (int(cameraX) + sw) / tw + 1
	endRow := (int(cameraY) + sh) / th + 1

	// Clamp to map bounds
	if startCol < 0 {
		startCol = 0
	}
	if startRow < 0 {
		startRow = 0
	}
	if endCol > tm.parsed.Width {
		endCol = tm.parsed.Width
	}
	if endRow > tm.parsed.Height {
		endRow = tm.parsed.Height
	}

	for _, layerName := range tm.renderOrder {
		grid := tm.parsed.Layers[layerName]
		for y := startRow; y < endRow; y++ {
			for x := startCol; x < endCol; x++ {
				gid := grid[y][x]
				if gid == 0 {
					continue // empty tile
				}

				srcRect := tm.tileset.TileRect(gid)
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(
					float64(x*tw)-cameraX,
					float64(y*th)-cameraY,
				)
				screen.DrawImage(
					tm.tileset.image.SubImage(srcRect).(*ebiten.Image),
					opts,
				)
			}
		}
	}
}

// GenerateColliders creates static physics bodies from the solid object layer.
func (tm *TileMap) GenerateColliders(world *physics.World) []*physics.Body {
	var bodies []*physics.Body

	for _, obj := range tm.parsed.SolidObjects {
		// TMX uses top-left origin; physics uses center origin
		centerX := obj.X + obj.Width/2
		centerY := obj.Y + obj.Height/2

		body := physics.NewBody(centerX, centerY, obj.Width, obj.Height, physics.Static)
		world.AddBody(body)
		bodies = append(bodies, body)
	}

	return bodies
}

// EntitySpawns returns the entity spawn points from the map.
func (tm *TileMap) EntitySpawns() []SpawnPoint {
	return tm.parsed.EntitySpawns
}

// MapWidthPx returns the total map width in pixels.
func (tm *TileMap) MapWidthPx() float64 {
	return float64(tm.parsed.Width * tm.parsed.TileWidth)
}

// MapHeightPx returns the total map height in pixels.
func (tm *TileMap) MapHeightPx() float64 {
	return float64(tm.parsed.Height * tm.parsed.TileHeight)
}
