#! /usr/bin/bash
set -e

source tmdb_key.sh
go run ./demo/. "${@}"