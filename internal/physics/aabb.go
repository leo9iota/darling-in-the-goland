package physics

import (
	"math"

	gm "github.com/leo9iota/darling-in-the-goland/internal/math"
)

// AABB is an Axis-Aligned Bounding Box defined by its min and max corners.
type AABB struct {
	Min, Max gm.Vec2
}

// NewAABB creates an AABB from a center position and dimensions (width, height).
// This matches Box2D's center-origin convention used in the Lua codebase.
func NewAABB(centerX, centerY, width, height float64) AABB {
	halfW := width / 2
	halfH := height / 2
	return AABB{
		Min: gm.Vec2{X: centerX - halfW, Y: centerY - halfH},
		Max: gm.Vec2{X: centerX + halfW, Y: centerY + halfH},
	}
}

// Overlaps returns true if this AABB overlaps with another (strict overlap, not touching).
func (a AABB) Overlaps(other AABB) bool {
	return a.Min.X < other.Max.X &&
		a.Max.X > other.Min.X &&
		a.Min.Y < other.Max.Y &&
		a.Max.Y > other.Min.Y
}

// ComputeMTV returns the Minimum Translation Vector needed to push this AABB out of
// the other. The returned vector points from other toward this AABB. If the AABBs don't
// overlap, a zero vector is returned.
func (a AABB) ComputeMTV(other AABB) gm.Vec2 {
	if !a.Overlaps(other) {
		return gm.Zero
	}

	// Compute overlap on each axis
	overlapLeft := a.Max.X - other.Min.X
	overlapRight := other.Max.X - a.Min.X
	overlapTop := a.Max.Y - other.Min.Y
	overlapBottom := other.Max.Y - a.Min.Y

	// Find smallest overlap on X axis
	var mtvX float64
	if overlapLeft < overlapRight {
		mtvX = -overlapLeft
	} else {
		mtvX = overlapRight
	}

	// Find smallest overlap on Y axis
	var mtvY float64
	if overlapTop < overlapBottom {
		mtvY = -overlapTop
	} else {
		mtvY = overlapBottom
	}

	// Push along the axis with smallest penetration
	if math.Abs(mtvX) < math.Abs(mtvY) {
		return gm.Vec2{X: mtvX, Y: 0}
	}
	return gm.Vec2{X: 0, Y: mtvY}
}
