#!/bin/bash

echo running completer
docker rm -f completer || true 
docker run --rm  -d -p 8081:8081 \
       -e DB_URL=inmem: \
       -e API_URL="http://$DOCKER_LOCALHOST:8080/r" \
       --name completer \
       fnproject/completer:latest


echo running completer ui
docker rm -f completerui || true 
docker run --name completerui --rm   -p3000:3000 -e API_URL=http://$DOCKER_LOCALHOST:8080 -d -e COMPLETER_BASE_URL=http://$DOCKER_LOCALHOST:8081 fnproject/completer:ui

echo open http://localhost:3000 for the completer UI