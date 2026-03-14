# Global Rules

## Terminal

- **Never use `cd`**: Always set the working directory via the `cwd` or equivalent parameter on each command. Do not prepend `cd <path> &&` to commands.
- **Avoid unnecessary commands**: Do not run build/test/lint commands unless the output is needed to proceed. Prefer letting the user run verification steps themselves (e.g. `task check`, `task run`).
- **Use Taskfile when available**: If a `taskfile.yml` exists in the project, use `task <name>` instead of raw commands (e.g. `task build` instead of `go build ./cmd/game`).
- **One command at a time**: Do not chain commands with `&&` or `;`. Run them as separate invocations so failures are isolated and recoverable.
- **Handle errors gracefully**: If a command fails, read the error output, diagnose, and fix — do not retry blindly or get stuck.
