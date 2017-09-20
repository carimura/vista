#!/bin/bash

DOCKER_LOCALHOST=docker.for.mac.localhost
echo running completer
docker rm -f completer || true 
docker run --rm  -d -p 8081:8081 \
       -e API_URL="http://$DOCKER_LOCALHOST:8080/r" \
       -e no_proxy=docker.for.mac.localhost \
       --name completer \
       fnproject/completer:latest


echo running completer ui
docker rm -f completerui || true 
docker run --name completerui --rm   -p3000:3000 -e API_URL=http://$DOCKER_LOCALHOST:8080 -d -e COMPLETER_BASE_URL=http://$DOCKER_LOCALHOST:8081 fnproject/completer:ui
