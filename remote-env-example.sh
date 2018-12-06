#!/usr/bin/env bash

set -ex

# remote distributed Fn
export FN_API_ENDPOINT=http://<your-api-service>:80
export FN_INVOKE_ENDPOINT=http://<your-lb-servie>:80/invoke

# all of this relates to Flow that you're suppose to run locally
export COMPLETER_BASE_URL=http://${DOCKER_LOCALHOST}:8081

# pubnub configuration
export PUBNUB_PUBLISH_KEY=".."
export PUBNUB_SUBSCRIBE_KEY=".."

# twitter configuration
export TWITTER_CONF_KEY="..."
export TWITTER_CONF_SECRET="..."
export TWITTER_TOKEN_KEY="..."
export TWITTER_TOKEN_SECRET="..."

# flickr configuration
export FLICKR_API_KEY="..."
export FLICKR_API_SECRET="..."

# slack configuration
export SLACK_API_TOKEN="...."

# ignore minio
export INSTALL_MINIO=0
export DOCKER_LOCALHOST=docker.for.mac.localhost
export MINIO_SERVER_URL=http://${DOCKER_LOCALHOST}:9000

# S3 access configuration
export STORAGE_ACCESS_KEY="DEMOACCESSKEY"
export STORAGE_SECRET_KEY="DEMOSECRETKEY"
export STORAGE_BUCKET="oracle-vista-out"
export S3_REGION="us-east-1"

# change this to deploy to a different app  other than "vista"
# export APP=vista

# set this to run vista in flow mode
export VISTA_MODE=flow
