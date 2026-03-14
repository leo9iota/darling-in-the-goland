package entity

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	gm "github.com/leo9iota/darling-in-the-goland/internal/math"

	"github.com/leo9iota/darling-in-the-goland/internal/animation"
	"github.com/leo9iota/darling-in-the-goland/internal/physics"
)

// Enemy constants from Enemy.lua.
const (
	enemySpeed       = 100.0
	enemyRageTrigger = 3
	enemyRageMult    = 3.0
	enemyDamage      = 1
	enemyOffsetY     = -8.0
	enemyGravityY    = 100.0 // constant downward push from Lua syncPhysics
)

// Enemy is a patrolling entity that damages the player.
type Enemy struct {
	Body *physics.Body

	anim      *animation.Controller
	xVel      float64 // current horizontal velocity (sign = direction)
	speed     float64
	speedMult float64

	rageCounter int

	spriteW float64
	spriteH float64
}

// NewEnemy creates a patrolling enemy at (x, y).
func NewEnemy(x, y float64, world *physics.World, player *Player) (*Enemy, error) {
	// Load walk animation (4 frames)
	walkFrames, err := loadFrames("assets/enemy/walk/%d.png", 4)
	if err != nil {
		return nil, fmt.Errorf("loading enemy walk: %w", err)
	}
	// Load run animation (4 frames)
	runFrames, err := loadFrames("assets/enemy/run/%d.png", 4)
	if err != nil {
		return nil, fmt.Errorf("loading enemy run: %w", err)
	}

	sw := float64(walkFrames[0].Bounds().Dx())
	sh := float64(walkFrames[0].Bounds().Dy())

	// Hitbox is 40% width, 75% height (matching Lua)
	hitW := sw * 0.4
	hitH := sh * 0.75

	body := physics.NewBody(x, y, hitW, hitH, physics.Dynamic)
	body.GravityScale = 0 // we handle gravity manually
	world.AddBody(body)

	anim := animation.NewController()
	anim.AddClip("walk", animation.NewClip(walkFrames, 0.1))
	anim.AddClip("run", animation.NewClip(runFrames, 0.07))

	e := &Enemy{
		Body:      body,
		anim:      anim,
		xVel:      enemySpeed,
		speed:     enemySpeed,
		speedMult: 1.0,
		spriteW:   sw,
		spriteH:   sh,
	}

	body.OnBeginContact = func(other *physics.Body, normal gm.Vec2) {
		if other == player.Body {
			player.TakeDamage(enemyDamage)
		}
		e.flipDirection()
		e.incrementRage()
	}

	body.UserData = e
	return e, nil
}

// Update advances enemy patrol and animation.
func (e *Enemy) Update(dt float64) {
	// Set velocity: patrol horizontally, constant downward push
	e.Body.Velocity.X = e.xVel * e.speedMult
	e.Body.Velocity.Y = enemyGravityY

	e.anim.Update(dt)
}

// Draw renders the enemy sprite, flipped based on direction.
func (e *Enemy) Draw(screen *ebiten.Image, camX, camY float64) {
	frame := e.anim.CurrentFrame()

	opts := &ebiten.DrawImageOptions{}

	// Flip if moving left
	if e.xVel < 0 {
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(e.spriteW, 0)
	}

	// Center sprite on body position with Y offset
	opts.GeoM.Translate(
		e.Body.Position.X-e.spriteW/2-camX,
		e.Body.Position.Y+enemyOffsetY-e.spriteH/2-camY,
	)

	screen.DrawImage(frame, opts)
}

func (e *Enemy) flipDirection() {
	e.xVel = -e.xVel
}

func (e *Enemy) incrementRage() {
	e.rageCounter++
	if e.rageCounter > enemyRageTrigger {
		e.anim.SetState("run")
		e.speedMult = enemyRageMult
		e.rageCounter = 0
	} else {
		e.anim.SetState("walk")
		e.speedMult = 1.0
	}
}
