# REFACTOR-001: Migration from Lua + LÖVE 2D to Go + Ebitengine

| Field       | Value                             |
| ----------- | --------------------------------- |
| **ID**      | REFACTOR-001                      |
| **Status**  | Draft                             |
| **Author**  | leo9iota                          |
| **Created** | 2026-03-14                        |
| **Source**  | Lua 5.4 / LÖVE 11.5 / STI / Box2D |
| **Target**  | Go 1.26 / Ebitengine v2           |

---

## 0: Master Checklist

### Phase 1: Project Bootstrap

- [x] Run `go mod init github.com/leo9iota/darling-in-the-goland`
- [x] Run `go get github.com/hajimehoshi/ebiten/v2`
- [x] Create `cmd/game/main.go`
  - [x] Define `Game` struct
  - [x] Implement `Update() error` (empty, return nil)
  - [x] Implement `Draw(screen *ebiten.Image)` (fill with background color)
  - [x] Implement `Layout(outsideW, outsideH int) (int, int)` (return 640×360 for 2× scaling)
  - [x] Call `ebiten.RunGame(&Game{})` in `main()`
- [x] Configure window properties
  - [x] Set title: `ebiten.SetWindowTitle("Darling in the GoLand")`
  - [x] Set size: `ebiten.SetWindowSize(1280, 720)`
  - [x] Disable vsync: `ebiten.SetVsyncEnabled(false)`
- [x] Handle ESC key to close window
- [x] Verify: `go run ./cmd/game` opens window, displays solid color, exits on ESC

### Phase 2: Asset Pipeline & Rendering Primitives

- [x] Create asset loading utilities
  - [x] Image loader: load PNG from disk → `*ebiten.Image`
  - [ ] Font loader: load TTF from disk → `text/v2` face
- [x] Create `internal/animation/sprite.go`
  - [x] Define `Animation` struct (frames, timer, rate, current frame index)
  - [x] `Update(dt)`: advance timer, cycle frame
  - [x] `CurrentFrame() *ebiten.Image`: return active frame
  - [x] `Reset()`: rewind to frame 0
- [x] Create `internal/core/background.go`
  - [x] Define `Layer` struct (image, parallax factor, x offset, y offset)
  - [x] `Load(paths, factors)`: load 4 layers, compute y from screen bottom
  - [x] `Update(cameraX)`: update each layer x = −cameraX × factor
  - [x] `Draw(screen)`: render layers back-to-front with tiling
- [ ] Load `public-pixel-font.ttf` as a reusable font face
- [x] Verify: background scrolls via arrow keys, an animated sprite plays on screen

### Phase 3: Custom Physics Engine

- [x] Create `internal/physics/aabb.go`
  - [x] Define `Vec2` struct (X, Y float64)
  - [x] Define `AABB` struct (Min, Max Vec2)
  - [x] `NewAABB(x, y, w, h)`: center-based constructor
  - [x] `Overlaps(other AABB) bool`: axis-aligned overlap test
  - [x] `Resolve(other AABB) Vec2`: minimum translation vector
- [x] Create `internal/physics/body.go`
  - [x] Define `BodyType` enum: `Static`, `Dynamic`, `Kinematic`
  - [x] Define `Body` struct
    - [x] Position, Velocity (Vec2)
    - [x] Width, Height (for AABB generation)
    - [x] BodyType, Mass, GravityScale
    - [x] FixedRotation, IsSensor flags
    - [x] OnBeginContact / OnEndContact callback fields
  - [x] `AABB() AABB`: compute AABB from position + dimensions
  - [x] `ApplyGravity(gravity Vec2, dt)`: add gravity to velocity
  - [x] `Integrate(dt)`: position += velocity × dt
- [x] Create `internal/physics/world.go`
  - [x] Define `World` struct (bodies slice, gravity Vec2)
  - [x] `NewWorld(gravX, gravY)`: constructor
  - [x] `AddBody(body)` / `RemoveBody(body)`
  - [x] `Update(dt)`: loop: apply gravity → integrate → detect collisions → resolve
  - [x] Track active contact pairs for begin/end detection
- [x] Create `internal/physics/collision.go`
  - [x] Define `Contact` struct (BodyA, BodyB, Normal Vec2, Depth float64)
  - [x] `DetectCollisions(bodies) []Contact`: brute-force broadphase
  - [x] `ResolveContact(contact)`: push dynamic bodies apart by MTV
  - [x] Skip resolution for sensor bodies (but still fire callbacks)
  - [ ] One-way platform logic: skip resolve if player velocity Y < 0
- [x] Write unit tests
  - [x] Test AABB overlap detection (overlapping, non-overlapping, edge cases)
  - [x] Test gravity integration (body falls at expected rate)
  - [x] Test collision resolution (body lands on floor, stops)
  - [x] Test sensor detection (callback fires, no position correction)
- [x] Verify: `go test ./internal/physics/...` all green

### Phase 4: Tilemap System

- [x] ~~Re-export Tiled maps to `.tmj` (JSON)~~ — parsed `.tmx` (XML) directly instead
  - [x] ~~Open each `.lua` map in Tiled~~ — used existing `.tmx` files in `assets/maps/tmx/`
  - [x] ~~Export as `.tmj` to `assets/maps/`~~ — no conversion needed
  - [x] Verify XML contains tile layers, object layers, tileset references
- [x] Create `internal/tilemap/loader.go`
  - [x] Define XML structs matching Tiled `.tmx` schema
  - [x] `LoadMap(path) (*TileMap, error)`: parse XML, resolve tileset paths
  - [x] Parse custom properties: `collidable`, entity `type`
- [x] Create `internal/tilemap/tileset.go`
  - [x] Load tileset source image as `*ebiten.Image`
  - [x] Compute tile rects: `TileRect(gid) image.Rectangle` (16×16 grid)
  - [x] Handle tileset first-GID offset
- [x] Create `internal/tilemap/tilemap.go`
  - [x] Define `TileMap` struct (layers, tilesets, width, height, tile size)
  - [x] `Draw(screen, cameraX, cameraY)`: render visible tile layers
  - [x] `GenerateColliders(world)`: create static bodies from `solid` layer
  - [x] `EntitySpawns() []SpawnPoint`: extract entity objects (type, x, y, w, h)
  - [x] Respect layer visibility (hide `solid` and `entity` layers)
- [x] Create `internal/core/camera.go`
  - [x] Define `Camera` struct (x, y, targetX, targetY, bounds, smoothing)
  - [x] `Follow(x, y)`: set target position, center on screen
  - [x] `Update(dt)`: exponential damped interpolation toward target
  - [x] Clamp to minX/maxX based on map width
  - [x] `SetBounds(mapWidth)`: compute maxX from map width and screen size
- [x] Compute `MapWidthPx = Width × TileWidth`
- [x] Verify: first level renders, camera follows a test point smoothly

### Phase 5: Player & Core Gameplay

- [ ] Create `internal/input/input.go`
  - [ ] `IsDown(keys ...ebiten.Key) bool`: any of the keys held
  - [ ] `IsJustPressed(key ebiten.Key) bool`: pressed this frame only
- [ ] Create `internal/entity/player.go`
  - [ ] Define `Player` struct
    - [ ] Physics body reference
    - [ ] Position, velocity, dimensions (20×60)
    - [ ] Health struct (current: 3, max: 3)
    - [ ] Coin count
    - [ ] Movement constants (maxSpeed: 200, accel: 1500, friction: 3500, gravity: 1500)
    - [ ] Jump constants (force: −650, doubleJumpMult: 0.75, coyoteTime: 0.25)
    - [ ] Animation state, direction, color tint
  - [ ] `Update(dt)` pipeline:
    - [ ] `untintRed(dt)`: restore color toward white
    - [ ] `respawn()`: reset if dead
    - [ ] `setAnimState()`: idle / run / air based on grounded + velocity
    - [ ] `setDirection()`: left / right based on xVelocity sign
    - [ ] `animate(dt)`: advance sprite animation
    - [ ] `decreaseCoyoteTimer(dt)`: count down when airborne
    - [ ] `syncPhysics()`: write velocity to physics body, read position back
    - [ ] `movement(dt)`: apply acceleration or friction based on input
    - [ ] `applyGravity(dt)`: add gravity when airborne
  - [ ] `Draw(screen, camera)`: draw current frame, flip if facing left, apply tint
  - [ ] `Jump()`: grounded/coyote → full jump; airborne + canDoubleJump → 0.75× jump
  - [ ] `TakeDamage(amount)`: reduce health, tint red, kill if health ≤ 0
  - [ ] `Kill()` / `Respawn()`: death flag, reset position + health
  - [ ] `IncrementCoinCount()`
- [ ] Wire collision callbacks
  - [ ] `OnBeginContact` → detect ground (normal pointing up), set grounded + reset coyote
  - [ ] `OnBeginContact` → detect ceiling (normal pointing down), zero Y velocity
  - [ ] `OnEndContact` → unground if leaving the ground collision
- [ ] Level transition logic: if `player.X > MapWidth - TileSize` → call `Map.Next()`
- [ ] Verify: full movement, jumping, double jump, coyote time, level transition

### Phase 6: Entities & Interactions

- [ ] Create `internal/entity/coin.go`
  - [ ] Define `Coin` struct (body, image, scaleX, randomTimeOffset, toRemove flag)
  - [ ] `NewCoin(x, y, world)`: create static sensor body
  - [ ] `Update(dt)`: spin animation: `scaleX = sin(time × 4 + offset)`
  - [ ] `Draw(screen, camera)`: draw coin with scaleX applied
  - [ ] On contact with player → set `toRemove = true`, increment player coins
  - [ ] Deferred removal: remove from world + slice after physics step
- [ ] Create `internal/entity/spike.go`
  - [ ] Define `Spike` struct (body, image, damage: 1)
  - [ ] `NewSpike(x, y, world)`: create static sensor body
  - [ ] On contact with player → `player.TakeDamage(1)`
  - [ ] `Draw(screen, camera)`: draw spike image
- [ ] Create `internal/entity/stone.go`
  - [ ] Define `Stone` struct (body, image, rotation)
  - [ ] `NewStone(x, y, world)`: create dynamic body, mass 25
  - [ ] `Update(dt)`: sync position + rotation from physics body
  - [ ] `Draw(screen, camera)`: draw with rotation applied
- [ ] Create `internal/entity/enemy.go`
  - [ ] Define `Enemy` struct (body, animations, state, speed, rageCounter, damage)
  - [ ] `NewEnemy(x, y, world)`: create dynamic body, fixed rotation, mass 25
  - [ ] `Update(dt)`: sync physics, animate, apply patrol velocity
  - [ ] `Draw(screen, camera)`: draw current frame, flip based on direction
  - [ ] `FlipDirection()`: negate xVel on any collision
  - [ ] `IncrementRage()`: count collisions, switch to run state (3× speed) after 3
  - [ ] On contact with player → `player.TakeDamage(1)`
- [ ] Create entity manager
  - [ ] Separate slices: `coins`, `enemies`, `spikes`, `stones`
  - [ ] `SpawnFromMap(spawns, world)`: instantiate by type string
  - [ ] `UpdateAll(dt)`: iterate each slice, update each entity
  - [ ] `DrawAll(screen, camera)`: iterate each slice, draw each entity
  - [ ] `RemoveMarked()`: remove entities flagged for deletion
  - [ ] `RemoveAll()`: destroy all bodies, clear slices (for level transitions)
  - [ ] `GetCounts() (coins, enemies, spikes, stones int)`: for debug overlay
- [ ] Wire entity spawning into tilemap entity layer
- [ ] Verify: all entities spawn, coins collectible, spikes damage, stones pushable, enemies patrol

### Phase 7: GUI & Polish

- [ ] Create `internal/gui/hud.go`
  - [ ] Load heart image + coin image
  - [ ] `Draw(screen, player)`: draw N hearts with spacing + shadow, coin icon + count
  - [ ] Shadow effect: draw image at +2px offset in black at 50% alpha
  - [ ] Use pixel font for coin count text
- [ ] Create `internal/gui/menu.go`
  - [ ] Define `Menu` struct (active, buttons with text + bounds + action)
  - [ ] 3 buttons: Resume (toggle menu), Settings (no-op print), Quit (os.Exit)
  - [ ] `Update()`: check mouse hover on buttons
  - [ ] `Draw(screen)`: semi-transparent overlay, button rects, hover highlight, text
  - [ ] `MousePressed(x, y)`: fire button action on click
  - [ ] `Toggle()`: flip active state
  - [ ] Block game updates when menu is active
- [ ] Create `internal/gui/debug.go`
  - [ ] Define `DebugGUI` struct (active, font, entity counts, metrics)
  - [ ] `Update(dt)`: read FPS (`ebiten.ActualFPS()`), compute frame time, memory stats
  - [ ] `UpdateEntityCounts(coins, enemies, spikes, stones)`
  - [ ] `Draw(screen)`: black background panel, yellow text, red for draw calls > 1000
  - [ ] `Toggle()`: F3 key binding
- [ ] Clean-up legacy files
  - [ ] Delete `main.lua`
  - [ ] Delete `conf.lua`
  - [ ] Delete `src/` directory (all Lua source)
  - [ ] Delete `modules/` directory (LuaRocks/STI)
- [ ] Update documentation
  - [ ] Rewrite `README.md` with Go build/run instructions
  - [ ] Update `docs/project/02-Tech-Stack.md` with Go + Ebitengine
- [ ] Final verification
  - [ ] `go vet ./...`: zero warnings
  - [ ] `go build ./cmd/game`: produces single executable
  - [ ] Full playthrough of all levels start to finish
  - [ ] All GUI elements functional (HUD, pause menu, debug overlay)

---

## 1: Motivation

The project name _"Darling in the GoLand"_ is a triple pun on the anime _Darling in the FranXX_, the Go programming language, and the JetBrains GoLand IDE. Currently the game is written in Lua + LÖVE 2D, which makes the name purely referential. **Rewriting the project in Go + Ebitengine** completes the joke and makes the name self-documenting.

Beyond the pun, Go provides meaningful technical advantages at this project's scale:

- **Static typing**: catches bugs at compile time rather than at runtime.
- **Single binary output**: no external runtime or framework installation required.
- **Native cross-compilation**: `GOOS=windows go build` and done.
- **Built-in testing & profiling**: `go test`, `go test -bench`, `pprof`.
- **No CGo dependency**: Ebitengine is pure Go (unlike raylib-go), simplifying builds.

---

## 2: Current System Inventory

A complete audit of every Lua source file and what it provides.

### 2.1 Core Systems

| File                      | LOC | Responsibilities                                                                 |
| ------------------------- | --- | -------------------------------------------------------------------------------- |
| `main.lua`                | 181 | Entry point, game loop (`load`/`update`/`draw`), input dispatch, collision relay |
| `conf.lua`                | 31  | LÖVE config: 1280×720 window, vsync off, title, LÖVE 11.5                        |
| `src/core/Map.lua`        | 101 | Level loading via STI + Box2D, tile layers, entity spawning, level transitions   |
| `src/core/Camera.lua`     | 204 | 2× scale, smooth follow (linear/damped), deadzone, bounds clamping               |
| `src/core/Background.lua` | 72  | 4-layer parallax scrolling, tiling, bottom-aligned rendering                     |
| `src/core/World.lua`      | 9   | Placeholder (physics world is created in `Map.lua` via `love.physics`)           |

### 2.2 Entities

| File                        | LOC | Responsibilities                                                                                                                                             |
| --------------------------- | --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `src/entities/Player.lua`   | 413 | Movement (accel/friction/gravity), double jump, coyote time, 3-state sprite animation (idle/run/air), health, damage tinting, collision handling, Box2D body |
| `src/entities/Enemy.lua`    | 151 | Patrol AI (walk/run), rage mechanic, direction flipping on collision, animated sprites, Box2D body                                                           |
| `src/entities/Coin.lua`     | 137 | Static collectible, spin animation via `scaleX`, sensor fixture, deferred removal                                                                            |
| `src/entities/Spike.lua`    | 68  | Static hazard, sensor fixture, deals damage on contact                                                                                                       |
| `src/entities/Stone.lua`    | 61  | Dynamic physics object, pushable by player, syncs rotation                                                                                                   |
| `src/entities/PowerUps.lua` | 9   | Placeholder: not yet implemented                                                                                                                             |

### 2.3 GUI

| File                   | LOC | Responsibilities                                              |
| ---------------------- | --- | ------------------------------------------------------------- |
| `src/gui/HUD.lua`      | 88  | Heart counter, coin counter, shadow rendering, pixel font     |
| `src/gui/Menu.lua`     | 149 | Pause overlay, 3 buttons (Resume/Settings/Quit), hover states |
| `src/gui/DebugGUI.lua` | 149 | FPS, frame time, memory, draw calls, entity counts            |

### 2.4 Assets

```
assets/
├── background/jungle/   (4 parallax layers)
├── enemy/run/           (4 frames)
├── enemy/walk/          (4 frames)
├── fonts/               (public-pixel-font.ttf)
├── images/              (icons, progress screenshots)
├── maps/                (Tiled .lua exports)
├── player/idle/         (4 frames)
├── player/run/          (6 frames)
└── world/               (coin.png, heart.png, spikes.png, stone.png)
```

### 2.5 External Dependencies

| Dependency         | Used For                 | Go Equivalent                  |
| ------------------ | ------------------------ | ------------------------------ |
| LÖVE 2D (Box2D)    | Physics world & bodies   | Custom AABB engine (hand-roll) |
| STI (Simple Tiled) | Tiled `.lua` map loading | Custom TMJ/TMX parser          |
| LuaRocks           | Package management       | Go modules (`go.mod`)          |

---

## 3: Target Architecture

### 3.1 Project Structure

```
darling-in-the-goland/
├── cmd/
│   └── game/
│       └── main.go              # Entry point
├── internal/
│   ├── core/
│   │   ├── game.go              # Game struct (implements ebiten.Game)
│   │   ├── camera.go            # Camera with smooth follow
│   │   └── background.go        # Parallax background
│   ├── physics/
│   │   ├── world.go             # Physics world (gravity, broadphase)
│   │   ├── body.go              # Rigid body (static, dynamic, kinematic)
│   │   ├── aabb.go              # AABB collision detection
│   │   └── collision.go         # Contact resolution, callbacks
│   ├── tilemap/
│   │   ├── loader.go            # Tiled JSON (.tmj) parser
│   │   ├── tilemap.go           # Tilemap struct, layer rendering
│   │   └── tileset.go           # Tileset image slicing
│   ├── entity/
│   │   ├── player.go            # Player entity
│   │   ├── enemy.go             # Enemy entity
│   │   ├── coin.go              # Coin collectible
│   │   ├── spike.go             # Spike hazard
│   │   └── stone.go             # Pushable stone
│   ├── gui/
│   │   ├── hud.go               # Hearts + coin counter
│   │   ├── menu.go              # Pause menu
│   │   └── debug.go             # Debug overlay (F3)
│   ├── animation/
│   │   └── sprite.go            # Frame-based sprite animation
│   └── input/
│       └── input.go             # Input abstraction (keyboard)
├── assets/                      # Unchanged: reuse all existing assets
├── docs/
├── go.mod
├── go.sum
└── README.md
```

### 3.2 Core Interface

Ebitengine requires implementing the `ebiten.Game` interface:

```go
type Game struct {
    // all game state lives here
}

func (g *Game) Update() error    // game logic (60 TPS by default)
func (g *Game) Draw(screen *ebiten.Image)  // rendering
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int)
```

This maps directly onto the current `love.update(dt)` / `love.draw()` pattern. **Key difference**: Ebitengine uses a fixed tick rate (default 60 TPS) rather than variable `dt`. Delta time is available via `1.0 / float64(ebiten.TPS())` when needed.

---

## 4: Migration Phases

Each phase produces a **playable checkpoint**: a version that compiles and runs, even if incomplete. No phase should exceed roughly one week of focused work.

---

### Phase 1: Project Bootstrap

**Goal**: Empty Ebitengine window, Go module initialized, project compiles and runs.

**Tasks**:

- [ ] Initialize Go module: `go mod init github.com/leo9iota/darling-in-the-goland`
- [ ] Add Ebitengine dependency: `go get github.com/hajimehoshi/ebiten/v2`
- [ ] Create `cmd/game/main.go` with minimal `ebiten.Game` implementation
- [ ] Configure window: 1280×720, title "Darling in the GoLand", pixel-art scaling
- [ ] Set `ebiten.SetVsyncEnabled(false)` to match current `conf.lua` (vsync = 0)
- [ ] Verify: window opens, solid background color, clean exit on ESC

**Acceptance Criteria**: `go run ./cmd/game` opens a 1280×720 window titled "Darling in the GoLand".

---

### Phase 2: Asset Pipeline & Rendering Primitives

**Goal**: Load and draw images, fonts, and sprite sheets. Parallax background renders.

**Tasks**:

- [ ] Create `internal/core/background.go`: port `Background.lua`
  - Load 4 jungle parallax layers as `*ebiten.Image`
  - Implement tiling with `ebiten.DrawImageOptions` + `GeoM.Translate`
  - Parallax factor per layer
- [ ] Create `internal/animation/sprite.go`: generic frame-based animation
  - Frame list, timer, rate, current frame tracking
  - Horizontal flip via `GeoM.Scale(-1, 1)`
- [ ] Set up asset embedding or loading from disk (`os.Open` / `ebitenutil.NewImageFromFile`)
- [ ] Load `public-pixel-font.ttf` via `text/v2` (Ebitengine font API)
- [ ] Verify: parallax background scrolls when simulating camera movement with arrow keys

**Acceptance Criteria**: Background layers render and scroll at different rates. At least one sprite draws and animates on screen.

---

### Phase 3: Custom Physics Engine

**Goal**: Replace Box2D with a hand-rolled AABB physics system suitable for a platformer.

**Tasks**:

- [ ] Create `internal/physics/aabb.go`
  - AABB struct: `Min`, `Max` as `Vec2`
  - Overlap test: `Overlaps(other AABB) bool`
  - Minimum Translation Vector (MTV): `Resolve(other AABB) Vec2`
- [ ] Create `internal/physics/body.go`
  - Body types: `Static`, `Dynamic`, `Kinematic`
  - Properties: position, velocity, gravity scale, fixed rotation, mass
  - Sensor flag (for coins, spikes: detect collision but don't resolve)
- [ ] Create `internal/physics/world.go`
  - Global gravity vector (0, 2000: matching current `love.physics.newWorld(0, 2000)`)
  - `Update(dt)`: apply gravity, integrate velocity, broadphase, narrowphase
  - Collision callbacks: `OnBeginContact`, `OnEndContact`
- [ ] Create `internal/physics/collision.go`
  - Contact struct: fixture pair, normal, penetration depth
  - Position correction (push bodies apart)
  - One-way platform support (only resolve if player is above)
- [ ] Verify: a test body falls under gravity and lands on a static platform

**Acceptance Criteria**: `go test ./internal/physics/...` passes. A dynamic body falls, collides with a static floor, and rests on it.

---

### Phase 4: Tilemap System

**Goal**: Load and render Tiled maps. Spawn entities from map object layers.

**Tasks**:

- [ ] Export existing Tiled maps to `.tmj` (JSON format) instead of `.lua`
  - Tiled supports this natively; the `.lua` exports are STI-specific
- [ ] Create `internal/tilemap/loader.go`
  - Parse `.tmj` JSON: layers, tilesets, objects
  - Map custom properties (`collidable`, entity `type`)
- [ ] Create `internal/tilemap/tileset.go`
  - Load tileset images, slice into sub-images by tile size (16×16)
- [ ] Create `internal/tilemap/tilemap.go`
  - Render tile layers using `ebiten.DrawImage` with source rects
  - Generate physics bodies from `solid` layer tiles
  - Parse `entity` layer objects → spawn positions (type + x/y/width/height)
  - Handle layer visibility (`solid` and `entity` layers hidden)
- [ ] Create `internal/core/camera.go`: port `Camera.lua`
  - Position, target, smoothing (damped spring), bounds, scale (2×)
  - `Follow(x, y)`, `Update(dt)`, apply via `GeoM` on draw
- [ ] Compute `MapWidth` from ground layer for camera bounds
- [ ] Verify: the first level renders at 2× scale, camera follows a moving point

**Acceptance Criteria**: `map-1` loads, tiles render, solid tiles generate collision bodies, camera pans across the level.

---

### Phase 5: Player & Core Gameplay

**Goal**: Playable character with full movement, animation, and collision.

**Tasks**:

- [ ] Create `internal/entity/player.go`: port `Player.lua`
  - Properties: position, velocity, dimensions (20×60), health (3), coins
  - Movement: acceleration (1500), friction (3500), max speed (200)
  - Gravity: 1500 units/s²
  - Jump: force −650, double jump (0.75× force), coyote time (0.25s)
  - Animation states: idle (4 frames), run (6 frames), air (4 frames)
  - Direction: flip sprite via `GeoM.Scale(-1, 1)`
  - Damage: red tint via `ColorScale`, untint over time
  - Death & respawn: reset position, restore health
- [ ] Create `internal/input/input.go`: keyboard abstraction
  - `IsDown(keys ...ebiten.Key) bool` for movement (A/D or ←/→)
  - `IsJustPressed(key ebiten.Key) bool` for jump (W or ↑), toggle (F3, ESC)
- [ ] Wire player physics body into the physics world
  - Dynamic body, fixed rotation
  - `beginContact` → grounding, wall detection
  - `endContact` → ungrounding
- [ ] Level transition: if `player.X > MapWidth - TileSize` → load next level
- [ ] Verify: player runs, jumps, double jumps, lands, and transitions to level 2

**Acceptance Criteria**: Player moves, jumps, animates correctly, collides with terrain, and falls off
ledges with appropriate coyote time.

---

### Phase 6: Entities & Interactions

**Goal**: All game entities ported and interacting with the player.

**Tasks**:

- [ ] Create `internal/entity/coin.go`: port `Coin.lua`
  - Static sensor body, spin animation via `scaleX = sin(time * 4 + offset)`
  - On contact with player: increment coin count, mark for removal
  - Deferred removal (remove after physics step, not during callback)
- [ ] Create `internal/entity/spike.go`: port `Spike.lua`
  - Static sensor body, deals 1 damage on player contact
- [ ] Create `internal/entity/stone.go`: port `Stone.lua`
  - Dynamic body, mass 25, pushed by player, syncs rotation
- [ ] Create `internal/entity/enemy.go`: port `Enemy.lua`
  - Dynamic body, patrol walk (speed 100), animated (walk 4f / run 4f)
  - Flip direction on any collision
  - Rage mechanic: after 3 collisions → run state (3× speed), resets
  - Deals 1 damage to player on contact
- [ ] Entity manager: `UpdateAll(dt)`, `DrawAll(screen)`, `RemoveAll()`, `GetCount()`
  - Separate slices for each entity type (matching `ActiveCoins`, `ActiveEnemies`, etc.)
- [ ] Map spawning: read entity layer objects, instantiate by `type` field
- [ ] Verify: all entity types spawn from the map, interact with player, and behave as expected

**Acceptance Criteria**: Coins collectible, spikes damage, stones pushable, enemies patrol and
deal damage. Entity counts visible in debug overlay.

---

### Phase 7: GUI & Polish

**Goal**: Full feature parity with the Lua version. Game is complete.

**Tasks**:

- [ ] Create `internal/gui/hud.go`: port `HUD.lua`
  - Heart counter: draw `N` hearts based on `player.Health.Current`
  - Coin counter: coin icon + count with pixel font
  - Shadow effect: draw at +2/+3px offset in black at 50% opacity
- [ ] Create `internal/gui/menu.go`: port `Menu.lua`
  - Pause overlay (ESC toggle), semi-transparent background
  - 3 buttons: Resume, Settings (no-op), Quit
  - Mouse hover detection, click handling
- [ ] Create `internal/gui/debug.go`: port `DebugGUI.lua`
  - FPS, frame time (ms), memory usage, draw calls (Ebitengine stats)
  - Entity counts (coins, enemies, spikes, stones, total)
  - Toggle with F3
  - Yellow text, warning color for high draw calls (>1000)
- [ ] Final clean-up
  - Remove all Lua source files (`main.lua`, `conf.lua`, `src/`, `modules/`)
  - Update `README.md` with Go build instructions
  - Update `docs/project/02-Tech-Stack.md`
- [ ] Verify: full playthrough of all levels from start to finish

**Acceptance Criteria**: Game is feature-complete, all GUI elements render, pause menu works,
debug overlay toggles. `go build ./cmd/game` produces a single executable.

---

## 5: Assets Reuse

All existing pixel art, fonts, and background images are **fully reusable**: they're standard PNG and TTF files. No Lua-specific asset formats are involved.

The only assets that need conversion are the **Tiled maps**: re-export from `.lua` format to `.tmj` (JSON) using the Tiled editor. The map geometry, layers, and object properties remain identical.

---

## 6: Lua ↔ Go Mapping Reference

A quick lookup table for porting patterns:

| Lua / LÖVE Pattern                                 | Go / Ebitengine Equivalent                                          |
| -------------------------------------------------- | ------------------------------------------------------------------- |
| `function love.load()`                             | `func main()` + `Game` struct init                                  |
| `function love.update(dt)`                         | `func (g *Game) Update() error`                                     |
| `function love.draw()`                             | `func (g *Game) Draw(screen *ebiten.Image)`                         |
| `love.graphics.newImage(path)`                     | `ebitenutil.NewImageFromFile(path)` or `ebiten.NewImageFromImage()` |
| `love.graphics.draw(img, x, y, r, sx, sy, ox, oy)` | `screen.DrawImage(img, &opts)` with `GeoM`                          |
| `love.graphics.push/pop + translate`               | `GeoM.Translate(x, y)` on draw options                              |
| `love.graphics.scale(2, 2)`                        | `GeoM.Scale(2, 2)` or `Layout()` return half-size                   |
| `love.graphics.setColor(r, g, b, a)`               | `opts.ColorScale.Scale(r, g, b, a)`                                 |
| `love.graphics.setFont(font)` / `print()`          | `text.Draw(screen, str, face, opts)`                                |
| `love.graphics.rectangle("fill", ...)`             | `ebitenutil.DrawRect(screen, x, y, w, h, color)`                    |
| `love.keyboard.isDown("d")`                        | `ebiten.IsKeyPressed(ebiten.KeyD)`                                  |
| `love.keypressed(key)`                             | `inpututil.IsKeyJustPressed(key)`                                   |
| `love.mouse.getPosition()`                         | `ebiten.CursorPosition()`                                           |
| `love.timer.getFPS()`                              | `ebiten.ActualFPS()`                                                |
| `love.timer.getTime()`                             | `time.Since(startTime).Seconds()`                                   |
| `setmetatable({}, Class)`                          | `type Struct struct { ... }` + methods                              |
| `table.insert(t, item)` / `#t`                     | `append(slice, item)` / `len(slice)`                                |
| `for i, v in ipairs(t)`                            | `for i, v := range slice`                                           |
| `math.min/max`                                     | `min(a, b)` / `max(a, b)` (built-in since Go 1.21)                  |
| `love.physics.newWorld(0, grav)`                   | `physics.NewWorld(0, grav)` (custom)                                |
| `love.physics.newBody(world, x, y, type)`          | `physics.NewBody(world, x, y, BodyDynamic)` (custom)                |
| `fixture:setSensor(true)`                          | `body.IsSensor = true` (custom)                                     |

---

## 7: Risk Assessment

| Risk                                      | Likelihood | Impact | Mitigation                                                |
| ----------------------------------------- | ---------- | ------ | --------------------------------------------------------- |
| Custom physics feels different than Box2D | Medium     | Medium | Tune constants iteratively; start with identical values   |
| Tiled JSON format has undocumented quirks | Low        | Low    | Only parse fields we actually need; ignore the rest       |
| Ebitengine draw performance at 2× scale   | Low        | Low    | Render to offscreen image at native res, scale up once    |
| Scope creep during rewrite                | Medium     | High   | Strict phase boundaries; each phase = playable checkpoint |
| Asset loading differences (paths, etc.)   | Low        | Low    | Use `os.Open` with relative paths from executable         |

---

## 8: Definition of Done

The migration is **complete** when:

1. `go run ./cmd/game` launches the game with identical gameplay to the Lua version.
2. All levels load and are playable from start to finish.
3. All entities behave identically (movement, damage, collection, physics).
4. All GUI elements render (HUD, pause menu, debug overlay).
5. No Lua source files remain in the repository.
6. `go build ./cmd/game` produces a single standalone executable.
7. `README.md` documents how to build and run the Go version.
