#!/usr/bin/env bash

set -ex

# on linux/bmc find this out with "docker inspect --type container -f '{{.NetworkSettings.Gateway}}' functions"
export DOCKER_LOCALHOST=docker.for.mac.localhost
export COMPLETER_BASE_URL=http://${DOCKER_LOCALHOST}:8081

fn update context registry `whoami`

export PUBNUB_PUBLISH_KEY=".."
export PUBNUB_SUBSCRIBE_KEY=".."

export TWITTER_CONF_KEY="..."
export TWITTER_CONF_SECRET="..."
export TWITTER_TOKEN_KEY="..."
export TWITTER_TOKEN_SECRET="..."

export FLICKR_API_KEY="..."
export FLICKR_API_SECRET="..."

export SLACK_API_TOKEN="...."

# Only change the following if you changed the defaults
export FUNC_SERVER_URL=http://${DOCKER_LOCALHOST}:8080
export MINIO_SERVER_URL=http://${DOCKER_LOCALHOST}:9000
export STORAGE_ACCESS_KEY="DEMOACCESSKEY"
export STORAGE_SECRET_KEY="DEMOSECRETKEY"
export STORAGE_BUCKET="oracle-vista-out"
export S3_REGION="us-east-1"

# change this to deploy to a different app  other than "vista"
# export APP=vista

# set this to run vista in flow mode
export VISTA_MODE=flow
