#! /usr/bin/bash
set -e

docker build -t video-in-be .
docker run --rm -v /nas/media:/nas/media -it -p 25004:25004 --name video-in-be video-in-be