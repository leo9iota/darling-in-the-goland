package gui

import (
	"fmt"
	"image/color"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Debug overlay constants — scaled for 640×360 internal resolution.
const (
	debugLineHeight = 11.0
	debugMargin     = 6.0
	debugPanelX     = 4.0
	debugPanelY     = 28.0 // below HUD hearts
)

var (
	debugTextColor = color.NRGBA{255, 255, 0, 191} // yellow 75% alpha
	debugBGColor   = color.NRGBA{0, 0, 0, 153}     // black 60% alpha
)

// Debug is a toggleable performance overlay.
type Debug struct {
	Active bool

	fps       float64
	frameTime float64
	memoryMB  float64

	coins   int
	enemies int
	spikes  int
	stones  int

	font *text.GoTextFace
}

// NewDebug creates a debug overlay (hidden by default).
func NewDebug() (*Debug, error) {
	fontSrc, err := LoadFont("assets/fonts/public-pixel-font.ttf")
	if err != nil {
		return nil, err
	}

	face := &text.GoTextFace{
		Source: fontSrc,
		Size:   FontSmall,
	}

	return &Debug{
		Active: false,
		font:   face,
	}, nil
}

// Toggle flips the debug overlay visibility.
func (d *Debug) Toggle() {
	d.Active = !d.Active
}

// Update refreshes performance metrics.
func (d *Debug) Update(dt float64) {
	if !d.Active {
		return
	}

	d.fps = ebiten.ActualFPS()
	d.frameTime = dt * 1000 // ms

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	d.memoryMB = float64(mem.Alloc) / 1024 / 1024
}

// UpdateEntityCounts sets the entity counts for display.
func (d *Debug) UpdateEntityCounts(coins, enemies, spikes, stones int) {
	d.coins = coins
	d.enemies = enemies
	d.spikes = spikes
	d.stones = stones
}

// Draw renders the debug panel.
func (d *Debug) Draw(screen *ebiten.Image) {
	if !d.Active {
		return
	}

	// Build lines
	total := d.coins + d.enemies + d.spikes + d.stones
	lines := []string{
		fmt.Sprintf("FPS: %.0f", d.fps),
		fmt.Sprintf("Frame: %.2fms", d.frameTime),
		fmt.Sprintf("Mem: %.1fMB", d.memoryMB),
		"",
		fmt.Sprintf("Coins: %d", d.coins),
		fmt.Sprintf("Enemies: %d", d.enemies),
		fmt.Sprintf("Spikes: %d", d.spikes),
		fmt.Sprintf("Stones: %d", d.stones),
		fmt.Sprintf("Total: %d", total),
	}

	// Measure widest line for panel width
	maxW := 0.0
	for _, l := range lines {
		if l == "" {
			continue
		}
		w, _ := MeasureText(l, d.font)
		if w > maxW {
			maxW = w
		}
	}

	pad := debugMargin
	panelW := maxW + pad*2
	panelH := float64(len(lines))*debugLineHeight + pad*2

	// Background panel
	vector.DrawFilledRect(screen, float32(debugPanelX), float32(debugPanelY), float32(panelW), float32(panelH), debugBGColor, true)

	y := debugPanelY + pad
	x := debugPanelX + pad

	for _, l := range lines {
		if l == "" {
			y += debugLineHeight * 0.5
			continue
		}
		d.drawLine(screen, l, x, y, debugTextColor)
		y += debugLineHeight
	}
}

func (d *Debug) drawLine(screen *ebiten.Image, str string, x, y float64, clr color.Color) {
	DrawText(screen, str, x, y, d.font, clr)
}
