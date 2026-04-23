#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env"

if ! command -v lazysql >/dev/null 2>&1; then
  echo "lazysql is required but was not found in PATH" >&2
  exit 1
fi

if [[ -f "${ENV_FILE}" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "${ENV_FILE}"
  set +a
fi

if [[ -z "${DATABASE_URL:-}" ]]; then
  echo "DATABASE_URL is required. Set it in .env or the environment." >&2
  exit 1
fi

if [[ $# -eq 0 ]]; then
  exec lazysql --read-only "${DATABASE_URL}"
fi

exec lazysql "$@" "${DATABASE_URL}"
