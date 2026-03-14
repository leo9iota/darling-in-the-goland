package physics

import (
	gm "github.com/leo9iota/darling-in-the-goland/internal/math"
)

// World manages the physics simulation.
type World struct {
	Gravity          gm.Vec2
	bodies           []*Body
	previousContacts map[contactPair]bool
	activeContacts   map[contactPair]bool
}

// contactPair identifies a unique pair of colliding bodies.
type contactPair struct {
	a, b *Body
}

func newContactPair(a, b *Body) contactPair {
	if uintptr(pointerOf(a)) > uintptr(pointerOf(b)) {
		a, b = b, a
	}
	return contactPair{a, b}
}

// NewWorld creates a physics world with the given gravity.
func NewWorld(gravX, gravY float64) *World {
	return &World{
		Gravity:          gm.Vec2{X: gravX, Y: gravY},
		previousContacts: make(map[contactPair]bool),
		activeContacts:   make(map[contactPair]bool),
	}
}

// AddBody adds a body to the world.
func (w *World) AddBody(b *Body) {
	w.bodies = append(w.bodies, b)
}

// RemoveBody removes a body from the world and cleans up contact tracking.
func (w *World) RemoveBody(b *Body) {
	for i, body := range w.bodies {
		if body == b {
			w.bodies = append(w.bodies[:i], w.bodies[i+1:]...)
			break
		}
	}

	for pair := range w.previousContacts {
		if pair.a == b || pair.b == b {
			delete(w.previousContacts, pair)
		}
	}
	for pair := range w.activeContacts {
		if pair.a == b || pair.b == b {
			delete(w.activeContacts, pair)
		}
	}
}

// Bodies returns the current list of bodies.
func (w *World) Bodies() []*Body {
	return w.bodies
}

// Update steps the physics simulation by dt seconds.
func (w *World) Update(dt float64) {
	// 1 & 2: Gravity + integration
	for _, b := range w.bodies {
		if b.Type != Dynamic {
			continue
		}
		b.Velocity.Y += w.Gravity.Y * b.GravityScale * dt
		b.Position.X += b.Velocity.X * dt
		b.Position.Y += b.Velocity.Y * dt
	}

	// 3: Detect collisions
	contacts := detectCollisions(w.bodies)

	// 4: Resolve contacts
	for _, c := range contacts {
		resolveContact(c)
	}

	// 5: Build active contact set and fire callbacks
	w.activeContacts = make(map[contactPair]bool, len(contacts))
	for _, c := range contacts {
		pair := newContactPair(c.BodyA, c.BodyB)
		w.activeContacts[pair] = true

		if !w.previousContacts[pair] {
			if c.BodyA.OnBeginContact != nil {
				c.BodyA.OnBeginContact(c.BodyB, c.Normal)
			}
			if c.BodyB.OnBeginContact != nil {
				c.BodyB.OnBeginContact(c.BodyA, c.Normal.Scale(-1))
			}
		}
	}

	for pair := range w.previousContacts {
		if !w.activeContacts[pair] {
			if pair.a.OnEndContact != nil {
				pair.a.OnEndContact(pair.b)
			}
			if pair.b.OnEndContact != nil {
				pair.b.OnEndContact(pair.a)
			}
		}
	}

	w.previousContacts = w.activeContacts
}
