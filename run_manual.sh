#! /usr/bin/bash

TAG="video-in-be:$(od -An -N4 -tu4 < /dev/urandom | tr -d ' ')"
echo "Building Docker image with tag: $TAG"
delete_docker_image() {
    echo "Deleting Docker image with tag: $TAG"
    docker rmi -f "$TAG" 2>/dev/null || true
}
docker build -t "$TAG" -f Dockerfile .
if [ $? -ne 0 ]; then
    echo "Docker build failed"
    exit 1
fi

source tmdb_key.sh
source fanart_key.sh

echo "Running Docker container with manual mode"
docker run --rm -it --name video-in-be-manual \
  -v /nas/media:/nas/media \
  --env VIDEOIN_TMDBKEY \
  --env VIDEOIN_FANARTAPIKEY \
  "$TAG" --mode=manual "${@}"
DOCKER_RUN_RESULT=$?

delete_docker_image

if [ $DOCKER_RUN_RESULT -ne 0 ]; then
    echo "Docker run failed with exit code $DOCKER_RUN_RESULT"
else 
    echo "Docker run completed successfully"
fi

exit $DOCKER_RUN_RESULT
