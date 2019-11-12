#!/usr/bin/env bash
# set -e

cd "$(dirname "$0")"

source ./.env/local.env

go test -v
