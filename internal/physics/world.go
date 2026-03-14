package physics

// World manages the physics simulation.
type World struct {
	Gravity          Vec2
	bodies           []*Body
	previousContacts map[contactPair]bool
	activeContacts   map[contactPair]bool
}

// contactPair identifies a unique pair of colliding bodies.
// Bodies are stored in pointer-order so (a,b) and (b,a) map to the same pair.
type contactPair struct {
	a, b *Body
}

func newContactPair(a, b *Body) contactPair {
	// Canonical ordering by pointer value to ensure (a,b) == (b,a)
	if uintptr(pointerOf(a)) > uintptr(pointerOf(b)) {
		a, b = b, a
	}
	return contactPair{a, b}
}

// NewWorld creates a physics world with the given gravity.
func NewWorld(gravX, gravY float64) *World {
	return &World{
		Gravity:          Vec2{gravX, gravY},
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

	// Clean up any contacts involving this body
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

// Bodies returns the current list of bodies (read-only use intended).
func (w *World) Bodies() []*Body {
	return w.bodies
}

// Update steps the physics simulation by dt seconds.
//
// The step order is:
//  1. Apply gravity to dynamic bodies
//  2. Integrate velocity into position
//  3. Detect collisions (N² broadphase)
//  4. Resolve contacts (push apart via MTV)
//  5. Fire begin/end contact callbacks
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

		// Fire OnBeginContact for new contacts
		if !w.previousContacts[pair] {
			if c.BodyA.OnBeginContact != nil {
				c.BodyA.OnBeginContact(c.BodyB, c.Normal)
			}
			if c.BodyB.OnBeginContact != nil {
				// Invert normal for body B's perspective
				c.BodyB.OnBeginContact(c.BodyA, c.Normal.Scale(-1))
			}
		}
	}

	// Fire OnEndContact for contacts that ended
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

	// Swap for next frame
	w.previousContacts = w.activeContacts
}
