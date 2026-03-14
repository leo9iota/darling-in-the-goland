# Tech Stack

## Core

| Technology                                    | Role               | Notes                                              |
| --------------------------------------------- | ------------------ | -------------------------------------------------- |
| [Go](https://go.dev/)                         | Language           | v1.26+, statically typed, single binary output     |
| [Ebitengine](https://ebitengine.org/)          | Game library       | v2, pure Go (no CGo), 2D rendering and input       |

## Game Systems

| System              | Implementation    | Description                                                   |
| ------------------- | ----------------- | ------------------------------------------------------------- |
| Physics             | Custom (in-house) | AABB collision detection, gravity, friction, sensor bodies     |
| Tilemap             | Custom (in-house) | Tiled `.tmj` (JSON) parser with layer rendering and colliders |
| Animation           | Custom (in-house) | Frame-based sprite animation with state machine               |
| Camera              | Custom (in-house) | Damped spring follow, deadzone, bounds clamping, 2× scale     |
| Parallax Background | Custom (in-house) | 4-layer depth-based scrolling with seamless tiling             |

## Tools

| Tool                                           | Role             | Notes                                                    |
| ---------------------------------------------- | ---------------- | -------------------------------------------------------- |
| [Tiled](https://www.mapeditor.org/)            | Level editor     | Exports `.tmj` (JSON), tile layers + object layers       |
| [Aseprite](https://www.aseprite.org/)          | Pixel art editor | Sprite sheets, animations, background layers             |
| [GoLand](https://www.jetbrains.com/go/)        | IDE              | JetBrains IDE for Go — also part of the project name pun |

## Assets

| Asset Type        | Format | Details                                       |
| ----------------- | ------ | --------------------------------------------- |
| Sprites           | PNG    | Pixel art, nearest-neighbor filtering          |
| Backgrounds       | PNG    | 4 jungle parallax layers                       |
| Fonts             | TTF    | Public Pixel Font for HUD and menus            |
| Maps              | TMJ    | Tiled JSON format with tile + entity layers    |

## Previous Stack (Pre-Migration)

> Prior to [REFACTOR-001](../tasks/todo/REFACTOR-001-Migration-to-Golang.md), the project used:

| Technology                                             | Role             |
| ------------------------------------------------------ | ---------------- |
| [Lua 5.4](https://www.lua.org/)                        | Language         |
| [LÖVE 2D 11.5](https://love2d.org/)                    | Game framework   |
| [Box2D](https://box2d.org/) (via LÖVE)                 | Physics engine   |
| [STI](https://github.com/karai17/Simple-Tiled-Implementation) | Tilemap loader |
| [LuaRocks](https://luarocks.org/)                      | Package manager  |
