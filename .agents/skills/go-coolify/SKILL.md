---
name: go-coolify
description: "Use go-coolify when writing Go code that calls the Coolify API or already imports github.com/psyb0t/go-coolify. The unofficial generated package coolify covers applications, deployments, databases, services, projects, servers and teams."
homepage: https://github.com/psyb0t/go-coolify
user-invocable: true
metadata:
  openclaw:
    emoji: "☁️"
    requires:
      bins:
        - go
      env:
        - COOLIFY_BASE_URL
        - COOLIFY_API_TOKEN
permissions:
  filesystem: "Guides edits to the Go source of the current project. The library reads and writes no files."
  shell: "The Go toolchain (`go get`, `go build`, `go test`) to add and compile the dependency."
  network: "Module installation uses your Go module proxy. Runtime requests go to the Coolify instance URL supplied by the user."
---

# go-coolify

Unofficial, typed Go client for the Coolify API, in package `coolify` at the module root.

Read [references/setup.md](references/setup.md) before adding the dependency or choosing operations. It covers authentication, base URLs, response handling and selected operations.

## Security and safety

- Read the user-supplied instance API URL from `COOLIFY_BASE_URL`, including `/api/v1`. There is no default host. Use HTTPS for remote instances and send credentials only to the intended instance.
- The spec's `bearerAuth` scheme uses `Authorization: Bearer <token>`. Read the token from `COOLIFY_API_TOKEN` in a request editor. Never hardcode, commit or log it, or print requests containing the authorization header.
- Creating, updating, deleting, deploying, restarting or stopping resources changes the user's infrastructure. Call these operations only when the user asked for the action.

## When to use

- Go code needs to inspect Coolify applications, deployments, databases, services, projects, servers or teams.
- Go code needs to change Coolify resources at the user's request.
- The project already imports `github.com/psyb0t/go-coolify`.

## When not to use

- The code is not Go. This module ships only a Go package.
- You need a client maintained by Coolify. This one is unofficial and not affiliated with Coolify.
- You want to change the generated client itself. Read the repository's `AGENTS.md` instead.

## Usage

Use the complete example in [references/setup.md](references/setup.md#usage). Construct `coolify.NewClientWithResponses` with the user's API URL and `coolify.WithRequestEditorFn`, then call `ListProjectsWithResponse(context.Background())` to list projects.

Check the error return for request, transport and decoding failures. HTTP errors may have a nil error, so also check `StatusCode()`, the expected `JSON...` field and `Body`. Check pointer fields before dereferencing. Close raw response bodies, and use `coolify.WithHTTPClient` with an explicit timeout.

Complete the change by compiling the caller and testing its request and response handling against a local HTTP test server. Do not contact the user's instance to validate an example without authorization.
