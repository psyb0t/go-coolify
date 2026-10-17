# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-17

- Added `coolify`, a REST client generated with oapi-codegen from Coolify's OpenAPI 3.1.0 document, spec version 0.1, at https://raw.githubusercontent.com/coollabsio/coolify/main/openapi.json. It covers all 202 paths and 291 operations, including applications, deployments, databases, services, projects, servers, teams and security resources.
- Unformatted number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` downloads the current spec, converts nullable type arrays and defines the missing 403 response component before regenerating. The spec is not committed.
- Named seven identical log-line query unions `LogTailLines` to expose generated helpers without changing their integer or `all` values. Description-only dash punctuation becomes commas.
- Added the `go-coolify` agent skill and plugin manifests for Claude Code and Codex.
