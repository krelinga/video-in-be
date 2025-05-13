#! /usr/bin/bash
set -e

docker build -t video-in-be .
docker run --rm -it -p 25004:25004 --name video-in-be-container video-in-be