# on linux/bmc find this out with "docker inspect --type container -f '{{.NetworkSettings.Gateway}}' functions"
export DOCKER_LOCALHOST=docker.for.mac.localhost

export FN_REGISTRY=<your docker id>

export PUBNUB_PUBLISH_KEY=X
export PUBNUB_SUBSCRIBE_KEY=X

export TWITTER_CONF_KEY=X
export TWITTER_CONF_SECRET=X
export TWITTER_TOKEN_KEY=X
export TWITTER_TOKEN_SECRET=X

export FLICKR_API_KEY=X
export FLICKR_API_SECRET=X

export SLACK_API_KEY=xoxb-....

# Only change the following if you changed the defaults
export FUNC_SERVER_URL=http://${DOCKER_LOCALHOST}:8080
export MINIO_SERVER_URL=http://${DOCKER_LOCALHOST}:9000
export STORAGE_ACCESS_KEY=DEMOACCESSKEY
export STORAGE_SECRET_KEY=DEMOSECRETKEY

# change this to deploy to a different app  other than "myapp"
# export APP=myapp

# set this to run vista in flow mode
# export VISTA_MODE=flow