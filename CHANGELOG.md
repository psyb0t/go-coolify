# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.6.0, 2026-10-29

- Breaking: `ServerSetting.IsLogdrainHighlightEnabled` and `ServerSetting.LogdrainHighlightProjectId` are removed. The spec no longer documents the Highlight log drain. Callers that read these fields must stop.
- Breaking: `IsPreview` is removed from `CreateEnvByServiceUuidJSONBody`, `UpdateEnvByServiceUuidJSONBody` and their `JSONRequestBody` aliases. Callers that set it must stop.
- Breaking: the `Data` element struct of `UpdateEnvsByServiceUuidJSONBody` and `UpdateEnvsByServiceUuidJSONRequestBody` loses its `IsPreview` field, so it is a different anonymous struct type. Code that names that type must change.

## v0.5.1, 2026-10-28

- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.5.0, 2026-10-27

- Breaking: `EnvironmentName` and `EnvironmentUuid` change from `string` to `*string` in the `JSONBody` and `JSONRequestBody` types of every create call for databases (`CreateDatabasePostgresqlJSONBody`, `CreateDatabaseMysqlJSONBody`, `CreateDatabaseRedisJSONBody` and the other engines) and applications (`CreatePublicApplicationJSONBody`, `CreateDockerfileApplicationJSONBody`, `CreateDockerimageApplicationJSONBody` and the other application types). Migration: assign a pointer, for example `&name`, and set only the one the spec requires.
- Breaking: `ListDatabasesResponse.JSON200` is now `*[]Database` and `GetDatabaseByUuidResponse.JSON200` is now `*Database`. Both were `*string`. The same change applies to the `GetJSON200` methods. Migration: read the typed value instead of parsing the string.
- Breaking: `ListResourcesResponse.JSON200` is now `*[]ListResources200JSONResponseBody_Item` instead of `*string`, with the matching `GetJSON200` change.
- Breaking: the 200 bodies of `GetS3StorageByUuidResponse`, `ListS3StoragesResponse` and `GetServerLogDrainsResponse` are different anonymous struct types, so code that names those types must change. The S3 structs gain `UsePathStyleEndpoint` and the log drain struct gains the Axiom domain field.
- Add the `Database` type, `DatabaseDatabaseType` and the constants `DatabaseDatabaseTypeStandalonePostgresql`, `DatabaseDatabaseTypeStandaloneMysql`, `DatabaseDatabaseTypeStandaloneMariadb`, `DatabaseDatabaseTypeStandaloneMongodb`, `DatabaseDatabaseTypeStandaloneRedis`, `DatabaseDatabaseTypeStandaloneKeydb`, `DatabaseDatabaseTypeStandaloneDragonfly`, `DatabaseDatabaseTypeStandaloneClickhouse` and `DatabaseDatabaseTypeStandaloneSqlite`.
- Add `UsePathStyleEndpoint` to `CreateS3StorageJSONBody`, `UpdateS3StorageByUuidJSONBody` and their `JSONRequestBody` aliases.
- Add `LogdrainAxiomDomain` to `ServerSetting`, `UpdateServerLogDrainsJSONBody` and `UpdateServerLogDrainsJSONRequestBody`.
- Add `GetHetznerServerTypesParams.Location`.
- Add the `N502` constant and `JSON502` fields with `GetJSON502` methods on the Hetzner list responses, `CreateHetznerServerResponse`, `CreateDigitaloceanServerResponse` and `CreateVultrServerResponse`.

## v0.4.2, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.

## v0.4.1, 2026-10-26

- Doc comments on the `Content` field of `UpdateStorageByApplicationUuidJSONBody`, `UpdateStorageByDatabaseUuidJSONBody` and `UpdateStorageByServiceUuidJSONBody` now say that changing the content needs the deploy permission. They no longer mention files outside the resource directory. No code changes.

## v0.4.0, 2026-10-24

- `DeleteServiceByUuidParams.DeleteFromCoolifyOnly` is removed. Callers that set it must stop setting it.
- The anonymous `JSON200` struct (and `GetJSON200`) of `GetServerProxyResponse`, `SaveServerProxyConfigurationResponse` and `UpdateServerProxyResponse` gets a new `ProxyErrorPage` field. Code that names these struct types by hand must be updated.
- `DeleteServiceByUuidResponse` gets `JSON422` and `GetJSON422`.
- `UpdateServerProxyJSONBody` and `UpdateServerProxyJSONRequestBody` get `ProxyErrorPage`, a custom HTML error page for unknown hosts and maintenance mode.
- `UpdateServerByUuidJSONBody` and `UpdateServerByUuidJSONRequestBody` get `BackupCompressionCpuPercentage`, `ServerTimezone` and `WildcardDomain`. The new type `UpdateServerByUuidJSONBodyBackupCompressionCpuPercentage` has the values `N25`, `N50`, `N75` and `N100`, plus a `Valid` method.
- `PatchServiceApplicationByServiceAndAppUuidJSONBody` and its request body get `Redirect`. The new type `PatchServiceApplicationByServiceAndAppUuidJSONBodyRedirect` has the values `Both`, `NonWww` and `Www`, plus a `Valid` method.
- The application create and update bodies (`CreatePublicApplicationJSONBody`, `CreatePrivateGithubAppApplicationJSONBody`, `UpdateApplicationByUuidJSONBody` and the matching `JSONRequestBody` types) get `DockerfileTargetBuild`, `IncludeCoolifyDefaultVariables`, `IsContainerLabelReadonlyEnabled` and `RollingUpdateDrainSeconds`.
- `UpdateDatabaseByUuidJSONBody` gets `IncludeCoolifyDefaultVariables` and `IsLogDrainEnabled`. `UpdateServiceByUuidJSONBody` gets `IncludeCoolifyDefaultVariables`.
- `UpdateDatabaseBackupJSONBody` and its request body get `DisableLocalBackup`, which deletes the local backup file after the S3 upload.

## v0.3.0, 2026-10-21

- `EnableApi` is removed from the spec, so `(*Client).EnableApi`, `(*ClientWithResponses).EnableApiWithResponse`, `EnableApiResponse`, `NewEnableApiRequest` and `ParseEnableApiResponse` are gone, along with the matching `ClientInterface` and `ClientWithResponsesInterface` methods.
- `UpdatePrivateKeyJSONBody.PrivateKey` and `UpdatePrivateKeyJSONRequestBody.PrivateKey` change from `string` to `*string` and are now optional.
- Server log drain responses (`GetServerLogDrainsResponse`) and `UpdateServerLogDrainsJSONBody` get CloudWatch fields: `IsLogdrainCloudwatchEnabled`, `LogdrainCloudwatchGroup`, `LogdrainCloudwatchRegion` and `LogdrainCloudwatchStreamPrefix`. `ServerSetting` gets the same four.
- `ListApplicationRollbackImagesResponse` images get a `Source` field, typed `ListApplicationRollbackImages200JSONResponseBodyImagesSource` with `Local` and `Registry` values.
- New Hostinger server methods: `CreateHostingerServer`, `GetHostingerCatalog`, `GetHostingerDataCenters`, `GetHostingerPostInstallScripts`, `GetHostingerSshKeys` and `GetHostingerTemplates`, each with a `WithResponse` variant. Cloud token types get `CreateCloudTokenJSONBodyProviderHostinger` and `ListCloudTokens200JSONResponseBodyProviderHostinger`.
- New server transfer methods: `CheckServerImport`, `CompleteServerTransfer`, `GetServerTransferStatus`, `ImportServerTransferBundle` and `MigrateServerBetweenInstances`, with `WithBody` and `WithResponse` variants where the endpoint takes a body.
- Application create and update bodies and `Application` get `HttpBasicAuthHashAlgorithm` (`bcrypt` or `argon2id`), `HttpBasicAuthBcryptCost`, `HttpBasicAuthArgon2idMemoryCost` and `HttpBasicAuthArgon2idTimeCost`.
- Application, service and database update bodies get `RedactAllEnvValuesInLogs`. `ApplicationSetting` and `Service` carry it too.
- `CreateDatabaseBackupJSONBody`, `UpdateDatabaseBackupJSONBody`, `VolumeBackupScheduleRequest` and `VolumeBackupScheduleResponse` get `S3StorageUuids`.
- The development image and `go.mod` move to Go 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by `make sec`. The setup reference in the skill names the new version.

## v0.2.0, 2026-10-19

- `CreateServerDestinationResponse.JSON422` and its `GetJSON422` method are removed. The spec no longer documents a 422 response for creating a server destination.
- `ApplicationSetting` gets a new `IncludeCoolifyDefaultVariables` field.
- `ApplicationSetting` and the application create, update and patch request bodies (`CreateDockerfileApplicationJSONBody`, `CreateDockerfileApplicationJSONRequestBody`, `CreateDockerimageApplicationJSONBody`, `CreateDockerimageApplicationJSONRequestBody`, `CreatePrivateDeployKeyApplicationJSONBody`, `CreatePrivateDeployKeyApplicationJSONRequestBody`, `CreatePrivateGithubAppApplicationJSONBody`, `CreatePrivateGithubAppApplicationJSONRequestBody`, `CreatePublicApplicationJSONBody`, `CreatePublicApplicationJSONRequestBody`, `PatchServiceApplicationByServiceAndAppUuidJSONBody`, `PatchServiceApplicationByServiceAndAppUuidJSONRequestBody`, `UpdateApplicationByUuidJSONBody`, `UpdateApplicationByUuidJSONRequestBody`) get a new `IsPermanentRedirectEnabled` field, to pick permanent 301/308 redirects instead of the default temporary 302/307 ones.
- Response parsers for several endpoints, for example `CreateDockerfileApplicationWithResponse` and `ExecuteScheduledTaskByServiceUuidWithResponse`, now recognize an HTTP 403 status explicitly. The spec still defines no body for it, so parsing just stops without an error.
- Generated doc comments for several methods now state the permission a token needs. For example, `DeleteApplicationByUuidWithResponse` needs the write and deploy permissions, and `SetServiceStorageBackupScheduleWithResponse` needs the deploy permission when the body sets `stop_during_backup`.

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
