package entity

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/leo9iota/darling-in-the-goland/internal/animation"
	"github.com/leo9iota/darling-in-the-goland/internal/input"
	gm "github.com/leo9iota/darling-in-the-goland/internal/math"
	"github.com/leo9iota/darling-in-the-goland/internal/physics"
)

// Player constants — all values ported from Player.lua.
const (
	playerWidth  = 20
	playerHeight = 60

	playerMaxSpeed     = 200.0
	playerAcceleration = 1500.0
	playerFriction     = 3500.0
	playerGravity      = 1500.0

	playerJumpForce      = -650.0
	playerDoubleJumpMult = 0.75
	playerCoyoteDuration = 0.25

	playerUntintSpeed = 3.0
)

// Player is the main character entity.
type Player struct {
	Body *physics.Body

	animations map[string]*animation.Animation
	animState  string
	direction  float64 // 1.0 = right, -1.0 = left

	isGrounded     bool
	canDoubleJump  bool
	coyoteTimer    float64
	groundContact  *physics.Body // which body we're standing on

	health    int
	maxHealth int
	coinCount int
	isAlive   bool

	startX, startY float64

	// Color tint (1,1,1 = white, 1,0,0 = red on damage)
	tintR, tintG, tintB float64
}

// NewPlayer creates a player at (x, y) and registers it with the physics world.
func NewPlayer(x, y float64, world *physics.World) (*Player, error) {
	body := physics.NewBody(x, y, playerWidth, playerHeight, physics.Dynamic)
	body.GravityScale = 0 // we handle gravity ourselves (matching Lua)
	world.AddBody(body)

	// Load animations
	idleFrames, err := loadFrames("assets/player/idle/zero-two-idle-%d.png", 4)
	if err != nil {
		return nil, fmt.Errorf("loading idle frames: %w", err)
	}
	runFrames, err := loadFrames("assets/player/run/zero-two-run-%d.png", 6)
	if err != nil {
		return nil, fmt.Errorf("loading run frames: %w", err)
	}
	// Air animation reuses idle frames (matching Lua)
	airFrames, err := loadFrames("assets/player/idle/zero-two-idle-%d.png", 4)
	if err != nil {
		return nil, fmt.Errorf("loading air frames: %w", err)
	}

	p := &Player{
		Body: body,
		animations: map[string]*animation.Animation{
			"idle": animation.NewAnimation(idleFrames, 0.1),
			"run":  animation.NewAnimation(runFrames, 0.1),
			"air":  animation.NewAnimation(airFrames, 0.1),
		},
		animState:   "idle",
		direction:   1.0,
		isGrounded:  false,
		canDoubleJump: true,
		health:      3,
		maxHealth:   3,
		isAlive:     true,
		startX:      x,
		startY:      y,
		tintR:       1, tintG: 1, tintB: 1,
	}

	// Register collision callbacks
	body.OnBeginContact = func(other *physics.Body, normal gm.Vec2) {
		if normal.Y < 0 {
			// Landed on ground (pushed up)
			p.land(other)
		} else if normal.Y > 0 {
			// Hit ceiling (pushed down)
			p.Body.Velocity.Y = 0
		}
	}
	body.OnEndContact = func(other *physics.Body) {
		if other == p.groundContact {
			p.isGrounded = false
		}
	}

	body.UserData = p

	return p, nil
}

// Update runs the full player update pipeline.
func (p *Player) Update(dt float64) {
	p.untintRed(dt)
	p.respawn()
	p.setAnimState()
	p.setDirection()
	p.movement(dt)
	p.applyGravity(dt)
	p.decreaseCoyoteTimer(dt)
	p.animations[p.animState].Update(dt)
	p.syncPhysics()
}

// Draw renders the player sprite at its position, offset by the camera.
func (p *Player) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	anim := p.animations[p.animState]
	frame := anim.CurrentFrame()
	fw := float64(frame.Bounds().Dx())
	fh := float64(frame.Bounds().Dy())

	opts := &ebiten.DrawImageOptions{}

	// Flip horizontally if facing left
	if p.direction < 0 {
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(fw, 0)
	}

	// Center sprite on physics body position
	opts.GeoM.Translate(
		p.Body.Position.X-fw/2-cameraX,
		p.Body.Position.Y-fh/2-cameraY,
	)

	// Apply color tint
	opts.ColorScale.Scale(float32(p.tintR), float32(p.tintG), float32(p.tintB), 1)

	screen.DrawImage(frame, opts)
}

// Jump handles jump input (called on key press, not hold).
func (p *Player) Jump() {
	if p.isGrounded || p.coyoteTimer > 0 {
		p.Body.Velocity.Y = playerJumpForce
		p.isGrounded = false
		p.coyoteTimer = 0
	} else if p.canDoubleJump {
		p.canDoubleJump = false
		p.Body.Velocity.Y = playerJumpForce * playerDoubleJumpMult
	}
}

// TakeDamage reduces health and tints the player red.
func (p *Player) TakeDamage(amount int) {
	p.tintRed()
	if p.health-amount > 0 {
		p.health -= amount
	} else {
		p.health = 0
		p.Kill()
	}
}

// Kill marks the player as dead.
func (p *Player) Kill() {
	p.isAlive = false
}

// Respawn resets the player if dead.
func (p *Player) Respawn() {
	if !p.isAlive {
		p.Body.Position.X = p.startX
		p.Body.Position.Y = p.startY
		p.Body.Velocity = gm.Vec2{}
		p.health = p.maxHealth
		p.isAlive = true
	}
}

// IncrementCoinCount adds one to the coin counter.
func (p *Player) IncrementCoinCount() {
	p.coinCount++
}

// CoinCount returns the current coin count.
func (p *Player) CoinCount() int {
	return p.coinCount
}

// Health returns the current health.
func (p *Player) Health() int {
	return p.health
}

// --- Internal methods ---

func (p *Player) land(other *physics.Body) {
	p.groundContact = other
	p.Body.Velocity.Y = 0
	p.isGrounded = true
	p.canDoubleJump = true
	p.coyoteTimer = playerCoyoteDuration
}

func (p *Player) untintRed(dt float64) {
	p.tintR = math.Min(p.tintR+playerUntintSpeed*dt, 1)
	p.tintG = math.Min(p.tintG+playerUntintSpeed*dt, 1)
	p.tintB = math.Min(p.tintB+playerUntintSpeed*dt, 1)
}

func (p *Player) tintRed() {
	p.tintG = 0
	p.tintB = 0
}

func (p *Player) respawn() {
	if !p.isAlive {
		p.Respawn()
	}
}

func (p *Player) setAnimState() {
	if !p.isGrounded {
		p.animState = "air"
	} else if p.Body.Velocity.X == 0 {
		p.animState = "idle"
	} else {
		p.animState = "run"
	}
}

func (p *Player) setDirection() {
	if p.Body.Velocity.X < 0 {
		p.direction = -1
	} else if p.Body.Velocity.X > 0 {
		p.direction = 1
	}
}

func (p *Player) movement(dt float64) {
	if input.IsDown(ebiten.KeyD, ebiten.KeyRight) {
		p.Body.Velocity.X = math.Min(p.Body.Velocity.X+playerAcceleration*dt, playerMaxSpeed)
	} else if input.IsDown(ebiten.KeyA, ebiten.KeyLeft) {
		p.Body.Velocity.X = math.Max(p.Body.Velocity.X-playerAcceleration*dt, -playerMaxSpeed)
	} else {
		p.applyFriction(dt)
	}
}

func (p *Player) applyFriction(dt float64) {
	if p.Body.Velocity.X > 0 {
		p.Body.Velocity.X = math.Max(p.Body.Velocity.X-playerFriction*dt, 0)
	} else if p.Body.Velocity.X < 0 {
		p.Body.Velocity.X = math.Min(p.Body.Velocity.X+playerFriction*dt, 0)
	}
}

func (p *Player) applyGravity(dt float64) {
	if !p.isGrounded {
		p.Body.Velocity.Y += playerGravity * dt
	}
}

func (p *Player) decreaseCoyoteTimer(dt float64) {
	if !p.isGrounded {
		p.coyoteTimer -= dt
	}
}

func (p *Player) syncPhysics() {
	// Velocity is already on the body; the world's Update() will integrate position.
	// We just ensure the body velocity is set from our calculated values.
	// (Position is read back from the body after world.Update in the game loop)
}

// loadFrames loads numbered PNG frames using a format pattern.
func loadFrames(pattern string, count int) ([]*ebiten.Image, error) {
	frames := make([]*ebiten.Image, 0, count)
	for i := 1; i <= count; i++ {
		path := fmt.Sprintf(pattern, i)
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("opening %s: %w", path, err)
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("decoding %s: %w", path, err)
		}
		frames = append(frames, ebiten.NewImageFromImage(img))
	}
	return frames, nil
}
