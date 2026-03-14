package entity

import (
	"fmt"
	"image"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	gm "github.com/leo9iota/darling-in-the-goland/internal/math"
	"github.com/leo9iota/darling-in-the-goland/internal/physics"
)

// Coin is a collectible entity that spins in place.
type Coin struct {
	Body     *physics.Body
	image    *ebiten.Image
	imgW     float64
	imgH     float64
	scaleX   float64
	time     float64
	offset   float64
	ToRemove bool
}

// NewCoin creates a spinning coin at (x, y).
// spawnX/Y are top-left from TMX; we center the body on the image.
func NewCoin(x, y float64, world *physics.World, player *Player) (*Coin, error) {
	img, w, h, err := loadEntityImage("assets/world/coin.png")
	if err != nil {
		return nil, err
	}

	cx := x + w/2
	cy := y + h/2

	body := physics.NewBody(cx, cy, w, h, physics.Static)
	body.IsSensor = true
	world.AddBody(body)

	c := &Coin{
		Body:   body,
		image:  img,
		imgW:   w,
		imgH:   h,
		scaleX: 1,
		offset: float64(int(x+y) % 100), // deterministic pseudo-random offset
	}

	body.OnBeginContact = func(other *physics.Body, _ gm.Vec2) {
		if other == player.Body {
			c.ToRemove = true
			player.IncrementCoinCount()
		}
	}

	body.UserData = c
	return c, nil
}

// Update advances the spin animation.
func (c *Coin) Update(dt float64) {
	c.time += dt
	c.scaleX = math.Sin(c.time*4 + c.offset)
}

// Draw renders the coin with horizontal scale (spin effect).
func (c *Coin) Draw(screen *ebiten.Image, camX, camY float64) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(c.scaleX, 1)
	// Center on body position
	opts.GeoM.Translate(-c.imgW*c.scaleX/2, -c.imgH/2)
	opts.GeoM.Translate(c.Body.Position.X-camX, c.Body.Position.Y-camY)
	screen.DrawImage(c.image, opts)
}

// loadEntityImage loads a PNG and returns the ebiten image + dimensions.
func loadEntityImage(path string) (*ebiten.Image, float64, float64, error) {
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
