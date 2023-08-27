#!/bin/bash

# Exit if error
set -e

docker compose build
docker image prune -f
docker compose up