package core

import (
	"math"
)

// Camera provides smooth-follow and viewport clamping for the game world.
type Camera struct {
	X, Y             float64
	targetX, targetY float64
	screenW, screenH float64
	minX, maxX       float64
	Smoothing        float64
}

// NewCamera creates a camera for the given screen dimensions.
func NewCamera(screenW, screenH float64) *Camera {
	return &Camera{
		screenW:   screenW,
		screenH:   screenH,
		Smoothing: 5.0, // higher = snappier follow
	}
}

// SetBounds configures the horizontal camera bounds based on the map width.
func (c *Camera) SetBounds(mapWidth float64) {
	c.minX = 0
	c.maxX = math.Max(0, mapWidth-c.screenW)
}

// Follow sets the target position. The camera will center on this point.
func (c *Camera) Follow(targetX, targetY float64) {
	c.targetX = targetX - c.screenW/2
	c.targetY = targetY - c.screenH/2
}

// Update smoothly interpolates the camera toward the target and clamps to bounds.
func (c *Camera) Update(dt float64) {
	// Damped follow
	t := 1 - math.Exp(-c.Smoothing*dt)
	c.X += (c.targetX - c.X) * t
	c.Y += (c.targetY - c.Y) * t

	// Clamp X
	if c.X < c.minX {
		c.X = c.minX
	}
	if c.X > c.maxX {
		c.X = c.maxX
	}

	// Clamp Y (don't scroll above the map)
	if c.Y < 0 {
		c.Y = 0
	}
}
