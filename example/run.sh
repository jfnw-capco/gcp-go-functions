#!/usr/bin/env bash

cd "$(dirname "$0")"

go build            # Builds
source .env.example # Set env
./example     # starts
