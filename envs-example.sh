#!/usr/bin/env bash

set -ex

# on linux/bmc find this out with "docker inspect --type container -f '{{.NetworkSettings.Gateway}}' functions"
export DOCKER_LOCALHOST=docker.for.mac.localhost

fn update context registry `whoami`

export PUBNUB_PUBLISH_KEY="pub-c-40a27c4c-2a77-42ac-9df7-027aac24f9b3"
export PUBNUB_SUBSCRIBE_KEY="sub-c-1a356ba8-8678-11e7-8979-5e3a640e5579"

export TWITTER_CONF_KEY=X
export TWITTER_CONF_SECRET=X
export TWITTER_TOKEN_KEY=X
export TWITTER_TOKEN_SECRET=X

export FLICKR_API_KEY="8fc6692d7f7390de4114ee4b84272d1a"
export FLICKR_API_SECRET="be8f0d8370f8ce2e"

export SLACK_API_TOKEN="...."

# Only change the following if you changed the defaults
export FUNC_SERVER_URL=http://${DOCKER_LOCALHOST}:8080
export MINIO_SERVER_URL=http://${DOCKER_LOCALHOST}:9000
export STORAGE_ACCESS_KEY=DEMOACCESSKEY
export STORAGE_SECRET_KEY=DEMOSECRETKEY

# change this to deploy to a different app  other than "vista"
# export APP=vista

# set this to run vista in flow mode
export VISTA_MODE=flow
