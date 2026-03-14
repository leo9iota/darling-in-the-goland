package physics

import (
	"testing"
)

func TestAABBOverlap(t *testing.T) {
	a := NewAABB(0, 0, 10, 10)  // -5..5 on both axes
	b := NewAABB(4, 0, 10, 10)  // -1..9 on X
	if !a.Overlaps(b) {
		t.Error("expected AABBs to overlap")
	}
}

func TestAABBNoOverlap(t *testing.T) {
	a := NewAABB(0, 0, 10, 10)
	b := NewAABB(20, 0, 10, 10) // 15..25 on X — no overlap
	if a.Overlaps(b) {
		t.Error("expected AABBs NOT to overlap")
	}
}

func TestAABBEdgeTouch(t *testing.T) {
	a := NewAABB(0, 0, 10, 10)  // -5..5
	b := NewAABB(10, 0, 10, 10) // 5..15 — edges touch, but not overlapping (strict <)
	if a.Overlaps(b) {
		t.Error("edge-touching AABBs should not count as overlap")
	}
}

func TestMTV(t *testing.T) {
	a := NewAABB(0, 0, 10, 10) // -5..5
	b := NewAABB(8, 0, 10, 10) // 3..13 on X → 2px overlap on X
	mtv := a.ComputeMTV(b)

	// Should resolve on X axis (smaller overlap than Y)
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

	// Dynamic body should have moved down
	if dynamic.Position.Y <= 0 {
		t.Errorf("expected dynamic body to fall, Y = %v", dynamic.Position.Y)
	}

	// Static body should not have moved
	if static.Position.Y != 100 {
		t.Errorf("expected static body to stay at Y=100, got %v", static.Position.Y)
	}
}

func TestCollisionResolution(t *testing.T) {
	w := NewWorld(0, 1000)

	// Player-like body falling onto a floor
	player := NewBody(0, 0, 10, 10, Dynamic)
	floor := NewBody(0, 10, 100, 5, Static) // just below the player

	w.AddBody(player)
	w.AddBody(floor)

	// Run several physics steps
	for i := 0; i < 120; i++ {
		w.Update(1.0 / 60.0)
	}

	// Player should be resting on the floor, not falling through
	floorTop := floor.Position.Y - floor.Height/2
	playerBottom := player.Position.Y + player.Height/2
	if playerBottom > floorTop+0.1 {
		t.Errorf("player fell through floor: playerBottom=%v, floorTop=%v", playerBottom, floorTop)
	}

	// Y velocity should be zeroed from collision resolution
	if player.Velocity.Y != 0 {
		t.Errorf("expected Y velocity to be 0 after landing, got %v", player.Velocity.Y)
	}
}

func TestSensor(t *testing.T) {
	w := NewWorld(0, 0) // no gravity for this test

	player := NewBody(0, 0, 10, 10, Dynamic)
	coin := NewBody(3, 0, 10, 10, Static)
	coin.IsSensor = true

	callbackFired := false
	coin.OnBeginContact = func(other *Body, normal Vec2) {
		if other == player {
			callbackFired = true
		}
	}

	w.AddBody(player)
	w.AddBody(coin)
	w.Update(1.0 / 60.0)

	// Callback should fire
	if !callbackFired {
		t.Error("expected sensor OnBeginContact to fire")
	}

	// Sensor should NOT push the player away — position should stay the same
	// (only gravity/velocity move dynamic bodies, no resolution for sensors)
	if player.Position.X != 0 || player.Position.Y != 0 {
		t.Errorf("sensor should not resolve: player position = %v, %v", player.Position.X, player.Position.Y)
	}
}

func TestBeginEndContact(t *testing.T) {
	w := NewWorld(0, 0)

	a := NewBody(0, 0, 10, 10, Dynamic)
	b := NewBody(5, 0, 10, 10, Static)
	b.IsSensor = true // use sensor so positions aren't adjusted

	beginCount := 0
	endCount := 0
	a.OnBeginContact = func(other *Body, normal Vec2) {
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
