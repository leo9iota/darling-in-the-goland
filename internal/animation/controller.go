package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Controller is a state machine that manages named animation clips.
// It only resets a clip when the state actually changes, preventing
// mid-animation resets from redundant SetState calls.
type Controller struct {
	clips   map[string]*Clip
	current string
}

// NewController creates an empty animation controller.
func NewController() *Controller {
	return &Controller{
		clips: make(map[string]*Clip),
	}
}

// AddClip registers a named clip. The first clip added becomes the default state.
func (c *Controller) AddClip(name string, clip *Clip) {
	c.clips[name] = clip
	if c.current == "" {
		c.current = name
	}
}

// SetState switches to a different clip. If the state is already active,
// this is a no-op — the animation continues uninterrupted.
func (c *Controller) SetState(name string) {
	if name == c.current {
		return
	}
	if clip, ok := c.clips[name]; ok {
		clip.Reset()
		c.current = name
	}
}

// State returns the name of the current animation state.
func (c *Controller) State() string {
	return c.current
}

// Update advances the current clip by dt seconds.
func (c *Controller) Update(dt float64) {
	if clip, ok := c.clips[c.current]; ok {
		clip.Update(dt)
	}
}

// CurrentFrame returns the active frame of the current clip.
func (c *Controller) CurrentFrame() *ebiten.Image {
	if clip, ok := c.clips[c.current]; ok {
		return clip.CurrentFrame()
	}
	return nil
}
