package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Animation holds a sequence of frames and cycles through them at a fixed rate.
type Animation struct {
	frames       []*ebiten.Image
	timer        float64
	rate         float64 // seconds per frame
	currentFrame int
}

// NewAnimation creates an animation from the given frames and rate (seconds per frame).
func NewAnimation(frames []*ebiten.Image, rate float64) *Animation {
	return &Animation{
		frames:       frames,
		rate:         rate,
		currentFrame: 0,
		timer:        0,
	}
}

// Update advances the animation timer by dt seconds and cycles frames.
func (a *Animation) Update(dt float64) {
	a.timer += dt
	for a.timer >= a.rate {
		a.timer -= a.rate
		a.currentFrame++
		if a.currentFrame >= len(a.frames) {
			a.currentFrame = 0
		}
	}
}

// CurrentFrame returns the active frame image.
func (a *Animation) CurrentFrame() *ebiten.Image {
	return a.frames[a.currentFrame]
}

// Reset rewinds the animation to the first frame.
func (a *Animation) Reset() {
	a.currentFrame = 0
	a.timer = 0
}

// FrameCount returns the total number of frames.
func (a *Animation) FrameCount() int {
	return len(a.frames)
}
