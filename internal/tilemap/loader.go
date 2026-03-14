package tilemap

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// --- TMX XML structs ---

// TMXMap is the root element of a .tmx file.
type TMXMap struct {
	XMLName    xml.Name         `xml:"map"`
	Width      int              `xml:"width,attr"`
	Height     int              `xml:"height,attr"`
	TileWidth  int              `xml:"tilewidth,attr"`
	TileHeight int              `xml:"tileheight,attr"`
	Tilesets   []TMXTileset     `xml:"tileset"`
	Layers     []TMXLayer       `xml:"layer"`
	ObjGroups  []TMXObjectGroup `xml:"objectgroup"`
}

// TMXTileset describes the tileset used by the map.
type TMXTileset struct {
	FirstGID   int       `xml:"firstgid,attr"`
	Name       string    `xml:"name,attr"`
	TileWidth  int       `xml:"tilewidth,attr"`
	TileHeight int       `xml:"tileheight,attr"`
	TileCount  int       `xml:"tilecount,attr"`
	Columns    int       `xml:"columns,attr"`
	Image      TMXImage  `xml:"image"`
}

// TMXImage is the tileset source image.
type TMXImage struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// TMXLayer is a tile layer with CSV-encoded data.
type TMXLayer struct {
	ID     int     `xml:"id,attr"`
	Name   string  `xml:"name,attr"`
	Width  int     `xml:"width,attr"`
	Height int     `xml:"height,attr"`
	Data   TMXData `xml:"data"`
}

// TMXData holds the CSV tile data.
type TMXData struct {
	Encoding string `xml:"encoding,attr"`
	Content  string `xml:",chardata"`
}

// TMXObjectGroup is a group of objects (solid colliders, entity spawns).
type TMXObjectGroup struct {
	ID      int         `xml:"id,attr"`
	Name    string      `xml:"name,attr"`
	Objects []TMXObject `xml:"object"`
}

// TMXObject is a single object in a group (rectangle, ellipse, or point).
type TMXObject struct {
	ID         int           `xml:"id,attr"`
	Name       string        `xml:"name,attr"`
	Type       string        `xml:"type,attr"`
	X          float64       `xml:"x,attr"`
	Y          float64       `xml:"y,attr"`
	Width      float64       `xml:"width,attr"`
	Height     float64       `xml:"height,attr"`
	Properties TMXProperties `xml:"properties"`
}

// TMXProperties wraps a list of properties.
type TMXProperties struct {
	Properties []TMXProperty `xml:"property"`
}

// TMXProperty is a custom property on an object.
type TMXProperty struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

// --- Parsed types ---

// SpawnPoint represents an entity spawn location extracted from the entity layer.
type SpawnPoint struct {
	Type          string
	X, Y          float64
	Width, Height float64
}

// ParsedMap holds the parsed data from a .tmx file.
type ParsedMap struct {
	Width, Height       int
	TileWidth, TileHeight int
	Layers              map[string][][]int // name → row-major tile grid
	SolidObjects        []TMXObject
	EntitySpawns        []SpawnPoint
	TilesetImagePath    string // resolved absolute path to tileset image
	TilesetFirstGID     int
	TilesetColumns      int
}

// LoadMap parses a .tmx file and returns a ParsedMap.
func LoadMap(path string) (*ParsedMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading tmx: %w", err)
	}

	var tmx TMXMap
	if err := xml.Unmarshal(data, &tmx); err != nil {
		return nil, fmt.Errorf("parsing tmx xml: %w", err)
	}

	if len(tmx.Tilesets) == 0 {
		return nil, fmt.Errorf("no tilesets found in %s", path)
	}

	ts := tmx.Tilesets[0]
	mapDir := filepath.Dir(path)
	tilesetPath := filepath.Join(mapDir, ts.Image.Source)

	parsed := &ParsedMap{
		Width:            tmx.Width,
		Height:           tmx.Height,
		TileWidth:        tmx.TileWidth,
		TileHeight:       tmx.TileHeight,
		Layers:           make(map[string][][]int),
		TilesetImagePath: tilesetPath,
		TilesetFirstGID:  ts.FirstGID,
		TilesetColumns:   ts.Columns,
	}

	// Parse tile layers
	for _, layer := range tmx.Layers {
		grid, err := parseCSV(layer.Data.Content, layer.Width, layer.Height)
		if err != nil {
			return nil, fmt.Errorf("parsing layer %q: %w", layer.Name, err)
		}
		parsed.Layers[layer.Name] = grid
	}

	// Parse object groups
	for _, group := range tmx.ObjGroups {
		switch group.Name {
		case "solid":
			for _, obj := range group.Objects {
				if hasProperty(obj, "collidable", "true") {
					parsed.SolidObjects = append(parsed.SolidObjects, obj)
				}
			}
		case "entity":
			for _, obj := range group.Objects {
				parsed.EntitySpawns = append(parsed.EntitySpawns, SpawnPoint{
					Type:   obj.Type,
					X:      obj.X,
					Y:      obj.Y,
					Width:  obj.Width,
					Height: obj.Height,
				})
			}
		}
	}

	return parsed, nil
}

// parseCSV converts Tiled's CSV tile data into a 2D grid.
func parseCSV(data string, width, height int) ([][]int, error) {
	data = strings.TrimSpace(data)
	grid := make([][]int, height)

	rows := strings.Split(data, "\n")
	for y := 0; y < height; y++ {
		row := strings.TrimSpace(rows[y])
		row = strings.TrimRight(row, ",")
		cols := strings.Split(row, ",")
		grid[y] = make([]int, width)
		for x := 0; x < width && x < len(cols); x++ {
			val, err := strconv.Atoi(strings.TrimSpace(cols[x]))
			if err != nil {
				return nil, fmt.Errorf("cell (%d,%d): %w", x, y, err)
			}
			grid[y][x] = val
		}
	}
	return grid, nil
}

// hasProperty checks if an object has a property with the given name and value.
func hasProperty(obj TMXObject, name, value string) bool {
	for _, p := range obj.Properties.Properties {
		if p.Name == name && p.Value == value {
			return true
		}
	}
	return false
}
