package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Clip holds a sequence of frames and cycles through them at a fixed rate.
type Clip struct {
	frames       []*ebiten.Image
	timer        float64
	rate         float64 // seconds per frame
	currentFrame int
	Looping      bool
}

// NewClip creates an animation clip from the given frames and rate (seconds per frame).
// Loops by default.
func NewClip(frames []*ebiten.Image, rate float64) *Clip {
	return &Clip{
		frames:       frames,
		rate:         rate,
		currentFrame: 0,
		timer:        0,
		Looping:      true,
	}
}

// Update advances the animation timer by dt seconds and cycles frames.
func (c *Clip) Update(dt float64) {
	c.timer += dt
	for c.timer >= c.rate {
		c.timer -= c.rate
		if c.Looping {
			c.currentFrame = (c.currentFrame + 1) % len(c.frames)
		} else if c.currentFrame < len(c.frames)-1 {
			c.currentFrame++
		}
	}
}

// CurrentFrame returns the active frame image.
func (c *Clip) CurrentFrame() *ebiten.Image {
	return c.frames[c.currentFrame]
}

// Reset rewinds the clip to the first frame.
func (c *Clip) Reset() {
	c.currentFrame = 0
	c.timer = 0
}

// IsFinished returns true if a non-looping clip has reached its last frame.
func (c *Clip) IsFinished() bool {
	return !c.Looping && c.currentFrame >= len(c.frames)-1
}

// FrameCount returns the total number of frames.
func (c *Clip) FrameCount() int {
	return len(c.frames)
}
