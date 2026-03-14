package physics

import (
	"testing"

	gm "github.com/leo9iota/darling-in-the-goland/internal/math"
)

func TestAABBOverlap(t *testing.T) {
	a := NewAABB(0, 0, 10, 10)
	b := NewAABB(4, 0, 10, 10)
	if !a.Overlaps(b) {
		t.Error("expected AABBs to overlap")
	}
}

func TestAABBNoOverlap(t *testing.T) {
	a := NewAABB(0, 0, 10, 10)
	b := NewAABB(20, 0, 10, 10)
	if a.Overlaps(b) {
		t.Error("expected AABBs NOT to overlap")
	}
}

func TestAABBEdgeTouch(t *testing.T) {
	a := NewAABB(0, 0, 10, 10)
	b := NewAABB(10, 0, 10, 10)
	if a.Overlaps(b) {
		t.Error("edge-touching AABBs should not count as overlap")
	}
}

func TestMTV(t *testing.T) {
	a := NewAABB(0, 0, 10, 10)
	b := NewAABB(8, 0, 10, 10)
	mtv := a.ComputeMTV(b)

	if mtv.X >= 0 {
		t.Errorf("expected negative X MTV, got %v", mtv.X)
	}
	if mtv.Y != 0 {
		t.Errorf("expected zero Y MTV, got %v", mtv.Y)
	}
}

func TestGravity(t *testing.T) {
	w := NewWorld(0, 1000)

	dynamic := NewBody(0, 0, 10, 10, Dynamic)
	static := NewBody(0, 100, 10, 10, Static)

	w.AddBody(dynamic)
	w.AddBody(static)

	dt := 1.0 / 60.0
	w.Update(dt)

	if dynamic.Position.Y <= 0 {
		t.Errorf("expected dynamic body to fall, Y = %v", dynamic.Position.Y)
	}
	if static.Position.Y != 100 {
		t.Errorf("expected static body to stay at Y=100, got %v", static.Position.Y)
	}
}

func TestCollisionResolution(t *testing.T) {
	w := NewWorld(0, 1000)

	player := NewBody(0, 0, 10, 10, Dynamic)
	floor := NewBody(0, 10, 100, 5, Static)

	w.AddBody(player)
	w.AddBody(floor)

	for i := 0; i < 120; i++ {
		w.Update(1.0 / 60.0)
	}

	floorTop := floor.Position.Y - floor.Height/2
	playerBottom := player.Position.Y + player.Height/2
	if playerBottom > floorTop+0.1 {
		t.Errorf("player fell through floor: playerBottom=%v, floorTop=%v", playerBottom, floorTop)
	}
	if player.Velocity.Y != 0 {
		t.Errorf("expected Y velocity to be 0 after landing, got %v", player.Velocity.Y)
	}
}

func TestSensor(t *testing.T) {
	w := NewWorld(0, 0)

	player := NewBody(0, 0, 10, 10, Dynamic)
	coin := NewBody(3, 0, 10, 10, Static)
	coin.IsSensor = true

	callbackFired := false
	coin.OnBeginContact = func(other *Body, normal gm.Vec2) {
		if other == player {
			callbackFired = true
		}
	}

	w.AddBody(player)
	w.AddBody(coin)
	w.Update(1.0 / 60.0)

	if !callbackFired {
		t.Error("expected sensor OnBeginContact to fire")
	}
	if player.Position.X != 0 || player.Position.Y != 0 {
		t.Errorf("sensor should not resolve: player position = %v, %v", player.Position.X, player.Position.Y)
	}
}

func TestBeginEndContact(t *testing.T) {
	w := NewWorld(0, 0)

	a := NewBody(0, 0, 10, 10, Dynamic)
	b := NewBody(5, 0, 10, 10, Static)
	b.IsSensor = true

	beginCount := 0
	endCount := 0
	a.OnBeginContact = func(other *Body, normal gm.Vec2) {
		beginCount++
	}
	a.OnEndContact = func(other *Body) {
		endCount++
	}

	w.AddBody(a)
	w.AddBody(b)

	// Frame 1: overlapping → begin contact
	w.Update(1.0 / 60.0)
	if beginCount != 1 {
		t.Errorf("expected 1 begin contact, got %d", beginCount)
	}

	// Frame 2: still overlapping → no new begin
	w.Update(1.0 / 60.0)
	if beginCount != 1 {
		t.Errorf("expected still 1 begin contact (no re-fire), got %d", beginCount)
	}

	// Move A far away → end contact
	a.Position.X = 100
	w.Update(1.0 / 60.0)
	if endCount != 1 {
		t.Errorf("expected 1 end contact, got %d", endCount)
	}
}
