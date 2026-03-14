# Git Commit Workflow

This guide establishes a rigorous commit discipline for IXIota development, emphasizing **atomic, logical commits** with **expressive messages**. Following these conventions ensures a readable, navigable history that makes debugging, code review, and collaboration significantly easier.

---

## Why Atomic Commits Matter

A well-crafted commit history is a developer's most powerful debugging tool. When each commit represents a single, coherent change:

- **Bisecting becomes trivial**: `git bisect` can pinpoint exactly where a bug was introduced.
- **Code review is faster**: Reviewers understand exactly what changed and why.
- **Reverts are safe**: Reverting one logical change doesn't break unrelated features.
- **Blame is useful**: `git blame` points to meaningful changes, not formatting fixes.

> 🎯 **Golden Rule**: If you find yourself using "and" in your commit message, you're probably committing too much at once. Split it into multiple commits.

---

## Conventional Commit Types

The Conventional Commits specification provides a lightweight structure for commit messages:

| Type         | Description                                                                                                          |
| ------------ | -------------------------------------------------------------------------------------------------------------------- |
| **feat**     | A new feature for the user (e.g., `feat(api): add payment intent endpoint`)                                          |
| **fix**      | A bug fix for the user (e.g., `fix(qr): resolve QR reference generation`)                                            |
| **docs**     | Documentation changes only (e.g., `docs: update API specification`)                                                  |
| **style**    | Code style only—formatting, whitespace, missing semicolons. **Not for UI changes.**                                  |
| **refactor** | Code change that neither fixes a bug nor adds a feature (e.g., `refactor(core): extract payment strategy interface`) |
| **perf**     | A change that improves performance (e.g., `perf(db): add index on transaction_id`)                                   |
| **test**     | Adding or correcting tests (e.g., `test: add unit tests for idempotency service`)                                    |
| **build**    | Changes to build system or dependencies (e.g., `build: add Liquibase dependency`)                                    |
| **ci**       | Changes to CI configuration (e.g., `ci: update GitHub Actions workflow`)                                             |
| **chore**    | Maintenance tasks that don't modify source or test files (e.g., `chore: update .gitignore`)                          |
| **revert**   | Reverting a previous commit (e.g., `revert: "feat(crypto): remove ETH integration"`)                                 |

---

## Commit Message Structure

Every commit message must follow this format:

```text
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Components

1. **Type**: One of the commit types above (lowercase)
2. **Scope**: The specific area affected (optional but recommended)
3. **Description**: Short, imperative summary (what this commit does)
4. **Body**: Detailed explanation (use when change is complex)
5. **Footer**: Breaking changes or issue references

### Examples

```text
feat(api): add payment intent creation endpoint

Implements POST /api/v1/payment_intents with support for
QR_BILL, TWINT, and CRYPTO payment methods.

Closes #42
```

```text
fix(db): resolve optimistic locking race condition

The PaymentIntent entity was vulnerable to concurrent updates
causing inconsistent state. Added @Version field for Hibernate
optimistic locking.

Fixes #128
```

```text
docs(spec): add payment state machine diagram

Updates 01-FinOps.md with state transition diagram showing
all possible PaymentIntent states and valid transitions.
```

---

## The Atomic Commit Workflow

Follow this step-by-step process for every commit:

### Step 1: Review Your Changes

Before staging, review what you've changed:

```bash
# See exactly what changed
git diff

# Or view staged changes
git diff --cached
```

### Step 2: Identify Logical Units

Ask yourself:

- ❓ "Does this change do two things?" → Split it
- ❓ "Can I describe this in one sentence?" → Good commit
- ❓ "Would I ever need to revert just this part?" → Key atomic test

### Step 3: Stage Incrementally

Stage only the files for one logical change:

```bash
# Good: Stage related files together
git add backend/src/main/java/.../PaymentIntent.java
git add backend/src/main/java/.../PaymentIntentRepository.java

# Bad: Stage everything at once
git add .
```

### Step 4: Write an Expressive Message

Use the **imperative mood** (describe what the commit _does_, not what you _did_):

| ❌ Bad                       | ✅ Good                                              |
| ---------------------------- | ---------------------------------------------------- |
| `feat: added user login`     | `feat(auth): add user login endpoint`                |
| `fix: fixed the button`      | `fix(ui): resolve submit button disabled state`      |
| `refactor: refactored stuff` | `refactor(core): extract payment strategy interface` |
| `chore: working on the api`  | `feat(api): add idempotency key validation`          |

### Step 5: Commit

```bash
git commit -m "feat(api): add payment intent creation endpoint"
```

---

## Branch Naming Convention

Commit discipline starts with branch naming. Use:

```text
<type>/<ticket-id>-<short-description>
```

Examples:

- `feature/42-payment-intent-api`
- `fix/128-concurrent-payment-race`
- `docs/01-update-state-machine`
- `refactor/extract-payment-strategy`

---

## Common Pitfalls to Avoid

| Pitfall                        | Problem                               | Solution                                  |
| ------------------------------ | ------------------------------------- | ----------------------------------------- |
| **Massive commits**            | Impossible to review or revert        | Split into logical units                  |
| **WIP commits**                | Pollutes history with incomplete work | Use `git rebase` to squash before merging |
| **Vague messages**             | Can't understand history later        | Always include scope and specific area    |
| **Mixed changes**              | Unrelated fixes in one commit         | Stage incrementally, commit separately    |
| **Imperative mood violations** | "Added" instead of "Add"              | Write as instructions to the codebase     |

---

## Quick Reference

```bash
# Check what you changed
git diff

# Stage a specific file
git add <path>

# Stage related files
git add <path1> <path2>

# Commit with message
git commit -m "type(scope): description"

# Amend last commit (if not pushed)
git commit --amend

# Interactive rebase to squash/fix commits
git rebase -i HEAD~n
```

---

## Summary

1. **One logical change per commit** — never mix features with fixes
2. **Write expressive messages** — scope + what + why (if needed)
3. **Use imperative mood** — "add" not "added"
4. **Review before staging** — understand what you're committing
5. **Keep it small** — if it's more than 500 lines, consider splitting
