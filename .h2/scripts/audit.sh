#!/usr/bin/env bash
# audit.sh — thin wrapper delegating to tlcli audit
exec "$(git rev-parse --show-toplevel)/bin/tlcli" audit "$@"
