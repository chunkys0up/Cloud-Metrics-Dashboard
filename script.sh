#!/bin/bash

set -e

redis-server --daemonize yes

cd server

go run main.go

# chmod +x ./script.sh