# aha-go CLI Roadmap

This document tracks planned enhancements for the aha-go CLI.

## Completed

- [x] Core API client for Aha.io
- [x] CLI commands: products (list, get)
- [x] CLI commands: features (list, get, create, update)
- [x] CLI commands: releases (list, get)
- [x] CLI commands: ideas (list, get)
- [x] CLI commands: goals (list, get)
- [x] CLI commands: initiatives (list, get)
- [x] CLI commands: epics (list, get)
- [x] CLI commands: requirements (list, get, create, update, delete)
- [x] CLI commands: canvas (create, list, get, update, export)
- [x] Canvas rendering (D2, Mermaid, SVG)
- [x] CLI commands: users (list, get, me)
- [x] CLI commands: comments (list, get, create, delete)
- [x] Output format support (--output json/yaml/table)
- [x] Shell completions (bash, zsh, fish, powershell)
- [x] Makefile with build, test, lint, install targets
- [x] GitHub Actions CI/CD workflows
- [x] CLI unit tests
- [x] Browser automation for template creation (go-rod)
- [x] Predefined canvas templates (capability-stack, maturity-model, opportunity-patton, feature-canvas)

## In Progress

### Feature-Release Query Support

Enable querying features by release date or name, supporting aha-studio's cached query capabilities.

| Item | Status | Description |
|------|--------|-------------|
| `ListFeaturesDetailed` | 🔄 In Progress | SDK method to list features with full details including release |
| Feature release_id in sync | 🔄 In Progress | Ensure FeatureMeta includes release reference for sync |
| CLI `--release-date` flag | 🔲 Planned | `aha feature list --release-date 2026-10-31` |
| CLI `--release-name` flag | 🔲 Planned | `aha feature list --release-name "Q4 2026"` |

## Future Considerations

- MCP server for AI assistant integration
- Interactive TUI mode
- Bulk operations (import/export)
- Webhook management
- Custom field support
- Tasks/To-dos (list, get, create, update, delete)
- Pages (list, get)

## Technical Debt

### Browser Automation: go-rod to chromedp Migration

**Status:** Planned (non-blocking)

The `browser` package uses [go-rod](https://github.com/go-rod/rod) for UI automation to create canvas templates (Aha.io has no API for template creation). However, go-rod appears unmaintained:

- Last meaningful code update: 2024
- Last release: v0.116.2 (July 2024)
- Recent commits: Logo/sponsor updates only (2025-2026)
- Dependency conflict: Requires pinned `fetchup v0.2.4`; newer versions have breaking API changes

**Current mitigation:** Pin `fetchup v0.2.4` in go.mod.

**Migration plan:**

| Phase | Description |
|-------|-------------|
| 1. Research | Evaluate chromedp API for aha-go use cases |
| 2. Wrapper | Create high-level wrapper around chromedp (rod's appeal was its abstractions) |
| 3. Migrate | Replace go-rod with chromedp wrapper in `browser` package |
| 4. Test | Verify template creation still works |
| 5. Remove | Remove go-rod and fetchup pin |

**Why chromedp:**

- Written in Go (go-first approach)
- Actively maintained by Google
- No external binary dependencies
- Uses same Chrome DevTools Protocol as go-rod

**Non-blocking:** This migration is technical debt cleanup, not a release blocker. The pinned dependency works correctly.
