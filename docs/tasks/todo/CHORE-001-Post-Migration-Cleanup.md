# CHORE-001: Post-Migration Cleanup

| Field       | Value                          |
| ----------- | ------------------------------ |
| **ID**      | CHORE-001                      |
| **Status**  | Todo                           |
| **Author**  | leo9iota                       |
| **Created** | 2026-03-14                     |
| **Depends** | REFACTOR-001 (Complete)        |

---

## Summary

Update project documentation to reflect the Go + Ebitengine stack and clean up git hygiene. Legacy Lua source is kept in `docs/code/` for reference.

---

## Tasks

### Documentation Update

- [ ] Rewrite `README.md`
  - [ ] Project description referencing Go + Ebitengine
  - [ ] Build prerequisites (Go 1.26+, Task)
  - [ ] Build & run instructions (`task run`, `task build:windows`, etc.)
  - [ ] Controls section (movement, jump, Escape, F3)
  - [ ] Project structure overview
- [ ] Update `docs/project/02-Tech-Stack.md`
  - [ ] Replace Lua/LÖVE references with Go/Ebitengine
  - [ ] Document custom physics engine (no Box2D dependency)
  - [ ] Document `text/v2` for font rendering
  - [ ] Document animation controller architecture

### Git Hygiene

- [ ] Add `game.exe` and build output to `.gitignore`
- [ ] Verify clean `git status` after all changes
