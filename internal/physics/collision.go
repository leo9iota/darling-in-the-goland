package physics

import (
	"math"
	"unsafe"

	gm "github.com/leo9iota/darling-in-the-goland/internal/math"
)

// pointerOf returns a uintptr for pointer-based canonical ordering.
func pointerOf(b *Body) unsafe.Pointer {
	return unsafe.Pointer(b)
}

// Contact represents a collision between two bodies.
type Contact struct {
	BodyA  *Body
	BodyB  *Body
	Normal gm.Vec2 // points from BodyB toward BodyA
	Depth  float64 // penetration depth
}

// detectCollisions performs N² broadphase AABB overlap detection.
// Skips static-static pairs since they can never move.
func detectCollisions(bodies []*Body) []Contact {
	var contacts []Contact

	for i := 0; i < len(bodies); i++ {
		for j := i + 1; j < len(bodies); j++ {
			a := bodies[i]
			b := bodies[j]

			// Skip static-static (neither can move)
			if a.Type == Static && b.Type == Static {
				continue
			}

			aBox := a.AABB()
			bBox := b.AABB()

			if !aBox.Overlaps(bBox) {
				continue
			}

			mtv := aBox.ComputeMTV(bBox)

			// Normalize MTV to get collision normal
			depth := math.Sqrt(mtv.X*mtv.X + mtv.Y*mtv.Y)
			var normal gm.Vec2
			if depth > 0 {
				normal = gm.Vec2{X: mtv.X / depth, Y: mtv.Y / depth}
			}

			contacts = append(contacts, Contact{
				BodyA:  a,
				BodyB:  b,
				Normal: normal,
				Depth:  depth,
			})
		}
	}

	return contacts
}

// resolveContact pushes dynamic bodies apart using the MTV.
// Sensors fire callbacks but are not resolved.
func resolveContact(c Contact) {
	// Don't resolve sensors (coins, spikes)
	if c.BodyA.IsSensor || c.BodyB.IsSensor {
		return
	}

	mtv := c.Normal.Scale(c.Depth)

	aDynamic := c.BodyA.Type == Dynamic
	bDynamic := c.BodyB.Type == Dynamic

	switch {
	case aDynamic && !bDynamic:
		// Push A out of B
		c.BodyA.Position = c.BodyA.Position.Add(mtv)
		// Zero velocity along collision normal
		zeroVelocityAlongNormal(c.BodyA, c.Normal)

	case !aDynamic && bDynamic:
		// Push B out of A
		c.BodyB.Position = c.BodyB.Position.Sub(mtv)
		// Zero velocity along collision normal (inverted)
		zeroVelocityAlongNormal(c.BodyB, c.Normal.Scale(-1))

	case aDynamic && bDynamic:
		// Split the push between both
		half := mtv.Scale(0.5)
		c.BodyA.Position = c.BodyA.Position.Add(half)
		c.BodyB.Position = c.BodyB.Position.Sub(half)
		zeroVelocityAlongNormal(c.BodyA, c.Normal)
		zeroVelocityAlongNormal(c.BodyB, c.Normal.Scale(-1))
	}
}

// zeroVelocityAlongNormal stops movement in the direction of the collision normal.
func zeroVelocityAlongNormal(b *Body, normal gm.Vec2) {
	if math.Abs(normal.Y) > math.Abs(normal.X) {
		// Vertical collision
		b.Velocity.Y = 0
	} else {
		// Horizontal collision
		b.Velocity.X = 0
	}
}
