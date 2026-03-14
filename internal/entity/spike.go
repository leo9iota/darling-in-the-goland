package entity

import (
	"github.com/hajimehoshi/ebiten/v2"

	gm "github.com/leo9iota/darling-in-the-goland/internal/math"
	"github.com/leo9iota/darling-in-the-goland/internal/physics"
)

// Spike is a hazard that damages the player on contact.
type Spike struct {
	Body  *physics.Body
	image *ebiten.Image
	imgW  float64
	imgH  float64
}

// NewSpike creates a spike at (x, y).
func NewSpike(x, y float64, world *physics.World, player *Player) (*Spike, error) {
	img, w, h, err := loadEntityImage("assets/world/spikes.png")
	if err != nil {
		return nil, err
	}

	cx := x + w/2
	cy := y + h/2

	body := physics.NewBody(cx, cy, w, h, physics.Static)
	body.IsSensor = true
	world.AddBody(body)

	s := &Spike{
		Body:  body,
		image: img,
		imgW:  w,
		imgH:  h,
	}

	body.OnBeginContact = func(other *physics.Body, _ gm.Vec2) {
		if other == player.Body {
			player.TakeDamage(1)
		}
	}

	body.UserData = s
	return s, nil
}

// Draw renders the spike image centered on its body.
func (s *Spike) Draw(screen *ebiten.Image, camX, camY float64) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(
		s.Body.Position.X-s.imgW/2-camX,
		s.Body.Position.Y-s.imgH/2-camY,
	)
	screen.DrawImage(s.image, opts)
}
