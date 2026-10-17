#!/bin/bash
# Converts nullable type arrays to OpenAPI 3.0 form for oapi-codegen.
# The preview-domain update references an undefined reusable 403 response;
# add a description-only definition without inventing its response body.
# Promote log-line query unions to a named component so oapi-codegen emits
# exported union helpers instead of unusable anonymous structs.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR
readonly SPEC_IN="${SPEC_DIR}/openapi.json"
readonly SPEC_OUT="${SPEC_DIR}/openapi.patched.json"

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg msg "$*" \
		'{time:$time,level:$level,msg:$msg}' >&2
}

main() {
	if [[ ! -f "$SPEC_IN" ]]; then
		log ERROR "source spec missing"
		exit 1
	fi

	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.openapi.patched.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp before this local goes out of scope
	trap "rm -f -- '${tmp}'" EXIT

	# Only the nullable two-type form is supported. Fail on new type unions
	# rather than silently discarding a possible wire representation.
	# shellcheck disable=SC2016 # $tails and $ref belong to jq, not the shell
	jq '
      walk(
        if type == "object" and (.type? | type == "array") then
          if (.type | length == 2 and index("null") != null) then
            .type = (.type | map(select(. != "null")) | .[0])
            | .nullable = true
          else error("unsupported type union") end
        else . end
      )
      | .components.responses["403"] //= {"description": "Forbidden."}
      | [.. | objects | select(.in? == "query" and .name? == "lines"
          and (.schema.oneOf? | type == "array")) | .schema] as $tails
      | if ($tails | length) > 0 then
          if all($tails[]; . == $tails[0]) then
            .components.schemas.LogTailLines //= $tails[0]
            | if .components.schemas.LogTailLines != $tails[0] then
                error("LogTailLines component collision")
              else walk(
                if type == "object" and .in? == "query"
                    and .name? == "lines"
                    and (.schema.oneOf? | type == "array") then
                  .schema = {"$ref": "#/components/schemas/LogTailLines"}
                else . end
              ) end
          else error("log-line schemas differ; review named union patch") end
        else . end
    ' "$SPEC_IN" >"$tmp"

	chmod 0644 "$tmp"
	mv -- "$tmp" "$SPEC_OUT"
	log INFO "patched nullable types, 403 response and log-line unions"
}

main "$@"
