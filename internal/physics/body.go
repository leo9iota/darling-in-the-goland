package physics

// BodyType defines how a body behaves in the physics simulation.
type BodyType int

const (
	// Static bodies don't move and aren't affected by gravity (floors, walls).
	Static BodyType = iota
	// Dynamic bodies are affected by gravity and respond to collisions (player, enemies).
	Dynamic
	// Kinematic bodies move but aren't affected by gravity (moving platforms).
	Kinematic
)

// Body is a rigid body in the physics world.
type Body struct {
	Position Vec2
	Velocity Vec2
	Width    float64
	Height   float64

	Type         BodyType
	GravityScale float64 // multiplier for world gravity (default 1.0)
	IsSensor     bool    // sensors detect overlap but don't resolve (coins, spikes)

	// Callbacks fired by the world during collision detection.
	OnBeginContact func(other *Body, normal Vec2)
	OnEndContact   func(other *Body)

	// UserData allows linking the body back to a game entity.
	UserData interface{}
}

// NewBody creates a body at (x, y) with the given dimensions and type.
func NewBody(x, y, width, height float64, bodyType BodyType) *Body {
	return &Body{
		Position:     Vec2{x, y},
		Width:        width,
		Height:       height,
		Type:         bodyType,
		GravityScale: 1.0,
	}
}

// AABB computes the bounding box for this body based on its position and dimensions.
func (b *Body) AABB() AABB {
	return NewAABB(b.Position.X, b.Position.Y, b.Width, b.Height)
}
