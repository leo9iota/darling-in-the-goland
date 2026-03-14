package entity

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/leo9iota/darling-in-the-goland/internal/physics"
)

// Stone is a pushable physics object.
type Stone struct {
	Body  *physics.Body
	image *ebiten.Image
	imgW  float64
	imgH  float64
}

// NewStone creates a pushable stone at (x, y).
func NewStone(x, y float64, world *physics.World) (*Stone, error) {
	img, w, h, err := loadEntityImage("assets/world/stone.png")
	if err != nil {
		return nil, err
	}

	cx := x + w/2
	cy := y + h/2

	body := physics.NewBody(cx, cy, w, h, physics.Dynamic)
	body.GravityScale = 1.0
	world.AddBody(body)

	s := &Stone{
		Body:  body,
		image: img,
		imgW:  w,
		imgH:  h,
	}

	body.UserData = s
	return s, nil
}

// Draw renders the stone centered on its body.
func (s *Stone) Draw(screen *ebiten.Image, camX, camY float64) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(
		s.Body.Position.X-s.imgW/2-camX,
		s.Body.Position.Y-s.imgH/2-camY,
	)
	screen.DrawImage(s.image, opts)
}
