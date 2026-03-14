# CHORE-001: Post-Migration Cleanup

| Field       | Value                          |
| ----------- | ------------------------------ |
| **ID**      | CHORE-001                      |
| **Status**  | Complete                       |
| **Author**  | leo9iota                       |
| **Created** | 2026-03-14                     |
| **Depends** | REFACTOR-001 (Complete)        |

---

## Summary

Update project documentation to reflect the Go + Ebitengine stack and clean up git hygiene. Legacy Lua source is kept for reference.

---

## Tasks

### Documentation Update

- [x] Rewrite `README.md`
  - [x] Project description referencing Go + Ebitengine
  - [x] Build prerequisites (Go 1.26+, Task)
  - [x] Build & run instructions (`task run`, `task build:windows`, etc.)
  - [x] Controls section (movement, jump, Escape, F3)
  - [x] Project structure overview
- [x] Update `docs/project/02-Tech-Stack.md`
  - [x] Replace Lua/LÖVE references with Go/Ebitengine
  - [x] Document custom physics engine (no Box2D dependency)
  - [x] Document `text/v2` for font rendering
  - [x] Document animation controller architecture

### Git Hygiene

- [x] Add `game.exe` and build output to `.gitignore`
- [x] Verify clean `git status` after all changes
