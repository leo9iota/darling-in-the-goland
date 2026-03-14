package entity

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/leo9iota/darling-in-the-goland/internal/physics"
	"github.com/leo9iota/darling-in-the-goland/internal/tilemap"
)

// Manager owns all entity slices and handles their lifecycle.
type Manager struct {
	Coins   []*Coin
	Spikes  []*Spike
	Stones  []*Stone
	Enemies []*Enemy
}

// NewManager creates an empty entity manager.
func NewManager() *Manager {
	return &Manager{}
}

// SpawnFromMap instantiates entities from TMX spawn points.
func (m *Manager) SpawnFromMap(spawns []tilemap.SpawnPoint, world *physics.World, player *Player) {
	for _, sp := range spawns {
		switch sp.Type {
		case "coin":
			c, err := NewCoin(sp.X, sp.Y, world, player)
			if err != nil {
				log.Printf("failed to create coin at (%.0f,%.0f): %v", sp.X, sp.Y, err)
				continue
			}
			m.Coins = append(m.Coins, c)

		case "spikes":
			s, err := NewSpike(sp.X, sp.Y, world, player)
			if err != nil {
				log.Printf("failed to create spike at (%.0f,%.0f): %v", sp.X, sp.Y, err)
				continue
			}
			m.Spikes = append(m.Spikes, s)

		case "stone":
			s, err := NewStone(sp.X, sp.Y, world)
			if err != nil {
				log.Printf("failed to create stone at (%.0f,%.0f): %v", sp.X, sp.Y, err)
				continue
			}
			m.Stones = append(m.Stones, s)

		case "enemy":
			e, err := NewEnemy(sp.X, sp.Y, world, player)
			if err != nil {
				log.Printf("failed to create enemy at (%.0f,%.0f): %v", sp.X, sp.Y, err)
				continue
			}
			m.Enemies = append(m.Enemies, e)
		}
	}

	log.Printf("Spawned %d coins, %d spikes, %d stones, %d enemies",
		len(m.Coins), len(m.Spikes), len(m.Stones), len(m.Enemies))
}

// UpdateAll updates entities that need per-frame logic.
func (m *Manager) UpdateAll(dt float64) {
	for _, c := range m.Coins {
		c.Update(dt)
	}
	for _, e := range m.Enemies {
		e.Update(dt)
	}
	// Spikes and stones don't need Update (static/physics-driven)
}

// DrawAll renders all entities.
func (m *Manager) DrawAll(screen *ebiten.Image, camX, camY float64) {
	for _, c := range m.Coins {
		c.Draw(screen, camX, camY)
	}
	for _, s := range m.Spikes {
		s.Draw(screen, camX, camY)
	}
	for _, s := range m.Stones {
		s.Draw(screen, camX, camY)
	}
	for _, e := range m.Enemies {
		e.Draw(screen, camX, camY)
	}
}

// RemoveMarked removes coins flagged for removal and destroys their bodies.
func (m *Manager) RemoveMarked(world *physics.World) {
	alive := m.Coins[:0]
	for _, c := range m.Coins {
		if c.ToRemove {
			world.RemoveBody(c.Body)
		} else {
			alive = append(alive, c)
		}
	}
	m.Coins = alive
}

// RemoveAll destroys all entity bodies and clears slices.
func (m *Manager) RemoveAll(world *physics.World) {
	for _, c := range m.Coins {
		world.RemoveBody(c.Body)
	}
	for _, s := range m.Spikes {
		world.RemoveBody(s.Body)
	}
	for _, s := range m.Stones {
		world.RemoveBody(s.Body)
	}
	for _, e := range m.Enemies {
		world.RemoveBody(e.Body)
	}
	m.Coins = nil
	m.Spikes = nil
	m.Stones = nil
	m.Enemies = nil
}

// GetCounts returns the current count of each entity type.
func (m *Manager) GetCounts() (coins, enemies, spikes, stones int) {
	return len(m.Coins), len(m.Enemies), len(m.Spikes), len(m.Stones)
}
