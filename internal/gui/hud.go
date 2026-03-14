package gui

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// HUD constants — scaled for 640×360 internal resolution.
const (
	hudHeartX       = 10.0
	hudHeartY       = 8.0
	hudScale        = 1.5
	hudShadowOffset = 1.0
	hudCoinShadow   = 1.0
	hudCoinMarginX  = 10.0 // margin from right edge
	hudHeartSpacing = 4.0  // gap between hearts
)

// HUD draws the heart counter and coin counter.
type HUD struct {
	heartImg *ebiten.Image
	heartW   float64
	heartH   float64

	coinImg *ebiten.Image
	coinW   float64
	coinH   float64

	font *text.GoTextFace

	screenW float64
}

// NewHUD creates the HUD, loading heart/coin images and the pixel font.
func NewHUD(screenW, screenH float64) (*HUD, error) {
	heartImg, heartW, heartH, err := loadGUIImage("assets/world/heart.png")
	if err != nil {
		return nil, fmt.Errorf("loading heart: %w", err)
	}

	coinImg, coinW, coinH, err := loadGUIImage("assets/world/coin.png")
	if err != nil {
		return nil, fmt.Errorf("loading coin: %w", err)
	}

	fontSrc, err := LoadFont("assets/fonts/public-pixel-font.ttf")
	if err != nil {
		return nil, fmt.Errorf("loading HUD font: %w", err)
	}

	face := &text.GoTextFace{
		Source: fontSrc,
		Size:   FontMedium,
	}

	return &HUD{
		heartImg: heartImg,
		heartW:   heartW,
		heartH:   heartH,
		coinImg:  coinImg,
		coinW:    coinW,
		coinH:    coinH,
		font:     face,
		screenW:  screenW,
	}, nil
}

// Draw renders the HUD (hearts + coin counter).
func (h *HUD) Draw(screen *ebiten.Image, health, maxHealth, coinCount int) {
	h.drawHearts(screen, health)
	h.drawCoinCounter(screen, coinCount)
}

func (h *HUD) drawHearts(screen *ebiten.Image, health int) {
	scaledW := h.heartW * hudScale

	for i := 0; i < health; i++ {
		x := hudHeartX + float64(i)*(scaledW+hudHeartSpacing)

		// Shadow
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(hudScale, hudScale)
		opts.GeoM.Translate(x+hudShadowOffset, hudHeartY+hudShadowOffset)
		opts.ColorScale.Scale(0, 0, 0, 0.5)
		screen.DrawImage(h.heartImg, opts)

		// Heart
		opts2 := &ebiten.DrawImageOptions{}
		opts2.GeoM.Scale(hudScale, hudScale)
		opts2.GeoM.Translate(x, hudHeartY)
		screen.DrawImage(h.heartImg, opts2)
	}
}

func (h *HUD) drawCoinCounter(screen *ebiten.Image, coinCount int) {
	// Measure text to compute total width
	countStr := fmt.Sprintf("%d", coinCount)
	textW, textH := MeasureText(countStr, h.font)

	scaledCoinW := h.coinW * hudScale
	scaledCoinH := h.coinH * hudScale
	gap := 4.0 // gap between coin icon and text

	// Total width = coin + gap + text
	totalW := scaledCoinW + gap + textW

	// Anchor from the right edge
	coinX := h.screenW - hudCoinMarginX - totalW

	// Coin shadow
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(hudScale, hudScale)
	opts.GeoM.Translate(coinX+hudCoinShadow, hudHeartY+hudCoinShadow)
	opts.ColorScale.Scale(0, 0, 0, 0.5)
	screen.DrawImage(h.coinImg, opts)

	// Coin
	opts2 := &ebiten.DrawImageOptions{}
	opts2.GeoM.Scale(hudScale, hudScale)
	opts2.GeoM.Translate(coinX, hudHeartY)
	screen.DrawImage(h.coinImg, opts2)

	// Text — vertically centered with coin
	textX := coinX + scaledCoinW + gap
	textY := hudHeartY + scaledCoinH/2 - textH/2

	DrawTextShadow(screen, countStr, textX, textY, h.font, hudCoinShadow, color.White)
}

// loadGUIImage loads a PNG and returns the image + dimensions.
func loadGUIImage(path string) (*ebiten.Image, float64, float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("decoding %s: %w", path, err)
	}

	eimg := ebiten.NewImageFromImage(img)
	w := float64(eimg.Bounds().Dx())
	h := float64(eimg.Bounds().Dy())
	return eimg, w, h, nil
}
