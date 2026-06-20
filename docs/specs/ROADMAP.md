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
