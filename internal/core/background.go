package core

import (
	"image"
	_ "image/png" // register PNG decoder
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// layer is a single parallax background layer.
type layer struct {
	image          *ebiten.Image
	parallaxFactor float64
	x              float64
	y              float64
}

// Background manages multi-layer parallax scrolling.
type Background struct {
	layers       []layer
	screenWidth  int
	screenHeight int
}

// NewBackground loads background layers from the given image paths with corresponding
// parallax factors. Layers are ordered from closest (index 0) to farthest.
func NewBackground(layerPaths []string, parallaxFactors []float64, screenWidth, screenHeight int) (*Background, error) {
	bg := &Background{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}

	for i, path := range layerPaths {
		img, err := loadImage(path)
		if err != nil {
			return nil, err
		}

		factor := 0.5
		if i < len(parallaxFactors) {
			factor = parallaxFactors[i]
		}

		// Align layer to the bottom of the screen
		imgHeight := img.Bounds().Dy()
		yOffset := float64(screenHeight) - float64(imgHeight)

		bg.layers = append(bg.layers, layer{
			image:          img,
			parallaxFactor: factor,
			x:              0,
			y:              yOffset,
		})
	}

	return bg, nil
}

// Update recalculates layer positions based on camera X position.
func (bg *Background) Update(cameraX float64) {
	for i := range bg.layers {
		bg.layers[i].x = -cameraX * bg.layers[i].parallaxFactor
	}
}

// Draw renders all background layers from farthest to closest with seamless tiling.
func (bg *Background) Draw(screen *ebiten.Image) {
	// Draw from farthest (last) to closest (first)
	for i := len(bg.layers) - 1; i >= 0; i-- {
		l := &bg.layers[i]
		img := l.image
		imgWidth := img.Bounds().Dx()

		// Calculate how many copies needed to cover the screen
		copies := bg.screenWidth/imgWidth + 2

		// Calculate starting X position for seamless tiling
		startX := math.Floor(math.Mod(l.x, float64(imgWidth))) - float64(imgWidth)

		for j := 0; j <= copies; j++ {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(startX+float64(j*imgWidth), l.y)
			screen.DrawImage(img, opts)
		}
	}
}

// loadImage loads a PNG file from disk as an *ebiten.Image.
func loadImage(path string) (*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}
