#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env"

if ! command -v lazysql >/dev/null 2>&1; then
  echo "lazysql is required but was not found in PATH" >&2
  exit 1
fi

if [[ -z "${DATABASE_URL:-}" ]] && [[ -f "${ENV_FILE}" ]]; then
  while IFS= read -r line; do
    line="${line#"${line%%[![:space:]]*}"}"
    [[ -z "${line}" || "${line:0:1}" == "#" ]] && continue

    if [[ "${line}" == export[[:space:]]* ]]; then
      line="${line#export }"
      line="${line#"${line%%[![:space:]]*}"}"
    fi

    if [[ "${line}" == DATABASE_URL=* ]]; then
      DATABASE_URL="${line#DATABASE_URL=}"
      if [[ "${DATABASE_URL}" == \"*\" ]] && [[ "${DATABASE_URL}" == *\" ]]; then
        DATABASE_URL="${DATABASE_URL:1:${#DATABASE_URL}-2}"
      elif [[ "${DATABASE_URL}" == \'*\' ]] && [[ "${DATABASE_URL}" == *\' ]]; then
        DATABASE_URL="${DATABASE_URL:1:${#DATABASE_URL}-2}"
      fi
      export DATABASE_URL
      break
    fi
  done < "${ENV_FILE}"
fi

if [[ -z "${DATABASE_URL:-}" ]]; then
  echo "DATABASE_URL is required. Set it in .env or the environment." >&2
  exit 1
fi

if [[ $# -eq 0 ]]; then
  exec lazysql --read-only "${DATABASE_URL}"
fi

exec lazysql "$@" "${DATABASE_URL}"
