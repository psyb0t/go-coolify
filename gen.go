// Package coolify is an unofficial Go client for the Coolify API,
// generated from the public OpenAPI document.
package coolify

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/patch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.patched.json
