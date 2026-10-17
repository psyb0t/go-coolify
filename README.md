# go-coolify

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-coolify.svg)](https://pkg.go.dev/github.com/psyb0t/go-coolify)
[![CI](https://github.com/psyb0t/go-coolify/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-coolify/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-coolify/badges/version.svg)](https://github.com/psyb0t/go-coolify/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-coolify/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-coolify/badges/importers.svg)](https://github.com/psyb0t/go-coolify/blob/badges/importers.md)

Unofficial Go client for the Coolify API, generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the [OpenAPI document Coolify publishes](https://raw.githubusercontent.com/coollabsio/coolify/main/openapi.json). I am not affiliated with Coolify.

The package is `coolify`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-coolify
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-coolify"
)

func main() {
	auth := coolify.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"Bearer "+os.Getenv("COOLIFY_API_TOKEN"),
			)

			return nil
		},
	)

	client, err := coolify.NewClientWithResponses(
		os.Getenv("COOLIFY_BASE_URL"),
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.ListProjectsWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles. Set `COOLIFY_BASE_URL` to your instance's API URL, including `/api/v1`, for example `https://coolify.example.com/api/v1`. The client always needs a user-supplied base URL.

The error return covers request, transport and decoding failures. HTTP errors can return a nil error; check `StatusCode()`, the appropriate `JSON...` field and `Body`. Check pointer fields before dereferencing. With raw client methods, close the response body. Pass `coolify.WithHTTPClient` with your own `*http.Client` to set a timeout.

## Auth

The spec's `bearerAuth` scheme uses HTTP Bearer authentication. Supply a token for your Coolify instance through `COOLIFY_API_TOKEN`. The request editor above sets `Authorization: Bearer <token>`; the client does not read environment variables or add authentication itself. Keep tokens out of source code and logs, and use HTTPS for remote instances.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `https://raw.githubusercontent.com/coollabsio/coolify/main/openapi.json` to `api/openapi.json`.
- `scripts/patch-spec.sh` writes `api/openapi.patched.json`. It rewrites OpenAPI 3.1 nullable type arrays into a single `type` plus `nullable: true`, and supplies the missing `components.responses.403` with the description `Forbidden.` so references resolve. It adds no response body schema.
- The patch also names seven identical log-line query unions `LogTailLines`, preserving their integer or `all` alternatives while exposing usable union helpers. It replaces dash punctuation in descriptions with commas.
- oapi-codegen runs with `config.yaml`. The config maps unformatted `number` fields to `float64` rather than the default `float32`, and generates all operations.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Agent integrations

The [skill](.agents/skills/go-coolify) teaches a coding agent how to use go-coolify when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-coolify@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-coolify@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The release workflow publishes the skill to ClawHub:

```bash
openclaw skills install @psyb0t/go-coolify
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Coolify's published API specification. The specification itself is not redistributed here. Its metadata declares no license or terms of service. Vendored dependencies retain their own licenses, including GPL licenses in development tooling.
