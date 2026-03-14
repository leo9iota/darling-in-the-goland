# FEAT-001: Gameplay Features & Enhancements

| Field       | Value                          |
| ----------- | ------------------------------ |
| **ID**      | FEAT-002                       |
| **Status**  | Todo                           |
| **Author**  | leo9iota                       |
| **Created** | 2026-03-14                     |
| **Depends** | REFACTOR-001 (Complete)        |

---

## Summary

New gameplay features deferred during the migration. These require design work beyond the original Lua port and represent genuine additions to the game.

---

## Tasks

### Level System

- [ ] Level transition logic
  - [ ] Detect player reaching end of map (`player.X > MapWidth - TileSize`)
  - [ ] Define level order (map-1, map-2, etc.)
  - [ ] `TileMap.Next()`: unload current level, load next map
  - [ ] Reset entities, camera, and player position on transition
  - [ ] Handle final level (show victory screen or loop)
- [ ] Level select / restart mechanics

### Physics Enhancements

- [ ] One-way platforms
  - [ ] Add `OneWay` flag to `Body`
  - [ ] Skip collision resolution when player velocity Y < 0 (moving up)
  - [ ] Allow drop-through with down+jump input
  - [ ] Mark specific tilemap tiles as one-way
- [ ] Spatial partitioning for collision broadphase
  - [ ] Grid-based or quadtree to replace N² brute-force
  - [ ] Benchmark before/after with entity-heavy levels

### New Enemies & Entities

- [ ] Flying enemy type (sine-wave patrol, ignores gravity)
- [ ] Moving platforms (kinematic body, waypoint path)
- [ ] Breakable blocks (hit from below to destroy)
- [ ] Checkpoints (save respawn position mid-level)
- [ ] Collectible power-ups (speed boost, extra jump, invincibility)

### Audio

- [ ] Background music (loop per level)
- [ ] Sound effects (jump, coin collect, damage, enemy rage)
- [ ] Volume control in Settings menu

### Visual Polish

- [ ] Particle effects (coin collect sparkle, damage flash, dust on land)
- [ ] Screen shake on damage
- [ ] Smooth camera zoom for boss areas
- [ ] Animated tile layers (water, lava)

### Performance

- [ ] Tile rendering culling (only draw visible tiles)
  - [ ] Compute visible tile range from camera position
  - [ ] Skip off-screen tiles in draw loop
- [ ] Object pooling for frequently created/destroyed entities
- [ ] Profile and optimize hot paths with `pprof`
