# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.2, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.
- The setup reference names Coolify Cloud's base URL, `https://app.coolify.io/api/v1`, for accounts on Coolify Cloud.

## v0.1.1, 2026-10-17

- Generated doc comments keep the API's own description text. The spec patch no longer rewrites dashes or wording in descriptions.
- `make sec` uses Semgrep 1.181.0. Its pinned requirements now hash every published file, so the development image also builds on arm64.

## v0.1.0, 2026-10-17

- Added `coolify`, a REST client generated with oapi-codegen from Coolify's OpenAPI 3.1.0 document, spec version 0.1, at https://raw.githubusercontent.com/coollabsio/coolify/main/openapi.json. It covers all 202 paths and 291 operations, including applications, deployments, databases, services, projects, servers, teams and security resources.
- Unformatted number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` downloads the current spec, converts nullable type arrays and defines the missing 403 response component before regenerating. The spec is not committed.
- Named seven identical log-line query unions `LogTailLines` to expose generated helpers without changing their integer or `all` values. Description-only dash punctuation becomes commas.
- Added the `go-coolify` agent skill and plugin manifests for Claude Code and Codex.
