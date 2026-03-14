package gui

import (
	"bytes"
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// LoadFont loads a TTF font file and returns a GoTextFaceSource.
func LoadFont(path string) (*text.GoTextFaceSource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading font %s: %w", path, err)
	}
	src, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parsing font %s: %w", path, err)
	}
	return src, nil
}

// DrawText draws a string at (x, y) with the given face and color.
func DrawText(screen *ebiten.Image, str string, x, y float64, face *text.GoTextFace, clr color.Color) {
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(x, y)
	opts.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, str, face, opts)
}

// DrawTextShadow draws text with a drop shadow.
// Shadow is drawn at +offset in black at 50% alpha, then the main text on top.
func DrawTextShadow(screen *ebiten.Image, str string, x, y float64, face *text.GoTextFace, offset float64, clr color.Color) {
	// Shadow
	shadowOpts := &text.DrawOptions{}
	shadowOpts.GeoM.Translate(x+offset, y+offset)
	shadowOpts.ColorScale.ScaleWithColor(color.NRGBA{0, 0, 0, 128})
	text.Draw(screen, str, face, shadowOpts)

	// Main text
	DrawText(screen, str, x, y, face, clr)
}

// MeasureText returns the width and height of a string rendered with the given face.
func MeasureText(str string, face *text.GoTextFace) (w, h float64) {
	w, h = text.Measure(str, face, 0)
	return w, h
}
