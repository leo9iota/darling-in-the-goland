# Darling in the GoLand

Classic platformer game inspired by [Super Mario](https://supermario-game.com/) in a retro pixel art style.

## About

**Darling in the GoLand** is a 2D side-scrolling platformer built with [Go](https://go.dev/) and the [Ebitengine](https://ebitengine.org/) game library. The player controls Zero Two through multi-level jungle environments, collecting coins, avoiding hazards, and defeating enemies.

### Origin of the Name

The name is a triple reference:

1. **Darling in the FranXX** (ダーリン・イン・ザ・フランキス) — the mecha anime
2. **Go** — the programming language the game is written in
3. **GoLand** — the JetBrains IDE for Go development

## Features

- **Platformer physics** — custom-built AABB physics engine with gravity, friction, and acceleration
- **Multi-level progression** — seamless level transitions via Tiled map system
- **Sprite animation** — frame-based animation system (idle, run, air states)
- **Parallax backgrounds** — 4-layer jungle background with depth-based scrolling
- **Entity system** — coins, spikes, pushable stones, and patrolling enemies with rage mechanic
- **Player mechanics** — double jump, coyote time, damage tinting, health and respawn
- **HUD** — heart counter, coin counter with pixel font and shadow effects
- **Pause menu** — resume, settings, and quit with mouse interaction
- **Debug overlay** — FPS, frame time, memory, draw calls, and entity counts (F3)

## Controls

| Key            | Action       |
| -------------- | ------------ |
| A / ←          | Move left    |
| D / →          | Move right   |
| W / ↑          | Jump         |
| W / ↑ (midair) | Double jump  |
| ESC            | Pause / Menu |
| F3             | Debug overlay |

## How to Run

```sh
git clone https://github.com/leo9iota/darling-in-the-goland.git && cd darling-in-the-goland
go run ./cmd/game
```

### Build

```sh
go build -o darling ./cmd/game
./darling
```

## Project Structure

```
darling-in-the-goland/
├── cmd/game/           Entry point
├── internal/
│   ├── core/           Game loop, camera, background
│   ├── physics/        Custom AABB physics engine
│   ├── tilemap/        Tiled .tmj map loader and renderer
│   ├── entity/         Player, enemies, coins, spikes, stones
│   ├── gui/            HUD, pause menu, debug overlay
│   ├── animation/      Frame-based sprite animation
│   └── input/          Keyboard input abstraction
├── assets/             Sprites, backgrounds, fonts, maps
└── docs/               Project documentation and task specs
```

## Documentation

- [Tech Stack](./02-Tech-Stack.md)
- [Migration Spec (REFACTOR-001)](../tasks/todo/REFACTOR-001-Migration-to-Golang.md)
