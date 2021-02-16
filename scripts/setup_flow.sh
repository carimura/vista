#!/bin/bash

set -ex

echo running flow-service
docker rm -f flow || true 
docker run --rm  -d -p 8081:8081 \
       -e DB_URL=inmem: \
       -e API_URL="${FN_INVOKE_ENDPOINT}" \
       --name flow \
       fnproject/flow:latest


echo running flow ui
docker rm -f flowui || true 
docker run --name flowui --rm   -p3000:3000 -e API_URL=${FN_API_ENDPOINT} -d -e COMPLETER_BASE_URL=${COMPLETER_BASE_URL} fnproject/flow:ui

echo open http://localhost:3000 for the flow UI
