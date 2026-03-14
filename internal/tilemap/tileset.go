package tilemap

import (
	"fmt"
	"image"
	_ "image/png" // register PNG decoder
	"os"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// Tileset holds the tileset image and metadata for computing tile source rectangles.
type Tileset struct {
	image    *ebiten.Image
	firstGID int
	tileW    int
	tileH    int
	columns  int
}

// LoadTileset loads a tileset image and stores its metadata.
func LoadTileset(imagePath string, firstGID, tileW, tileH, columns int) (*Tileset, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("opening tileset image %s: %w", imagePath, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decoding tileset image: %w", err)
	}

	return &Tileset{
		image:    ebiten.NewImageFromImage(img),
		firstGID: firstGID,
		tileW:    tileW,
		tileH:    tileH,
		columns:  columns,
	}, nil
}

// TileRect returns the source rectangle in the tileset image for the given GID.
// Returns a zero rectangle if the GID is invalid.
func (ts *Tileset) TileRect(gid int) image.Rectangle {
	if gid < ts.firstGID {
		return image.Rectangle{}
	}
	idx := gid - ts.firstGID
	col := idx % ts.columns
	row := idx / ts.columns
	x := col * ts.tileW
	y := row * ts.tileH
	return image.Rect(x, y, x+ts.tileW, y+ts.tileH)
}
