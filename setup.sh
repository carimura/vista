#!/usr/bin/env bash

set -ex

requiredEnv=$(grep -p "^export" < $(dirname $0)/envs-example.sh | sed 's/export \([^=]*\)=.*/\1/')

for envVar in $requiredEnv; do
    if [ -z $(printenv $envVar)  ] ; then
	echo "missing environment variable  \"$envVar\" , have you followed the setup instructions?"
        echo "required variables are $requiredEnv"
	exit 1
    fi
done 

APP=${APP:-vista}
export APP

if [ $VISTA_MODE = flow ] ; then
   echo "Running in flow mode , unset VISTA_MODE to run async  and rerun setup to deploy async "
else
   echo "Running in async mode , set VISTA_MODE to 'flow'  and rerun setup to deploy in flow mode "
   export VISTA_MODE="async"
fi

DOCKER_LOCALHOST=${DOCKER_LOCALHOST:-docker.for.mac.localhost}

echo "Setting up app:  $APP with docker localhost $DOCKER_LOCALHOST"

# we need to have a publish function in place before starting minio, so deploy first
# time="2017-08-30T16:21:13Z" level=error msg="Initializing object layer failed" cause="Unable to initialize event notification. Unexpected response from webhook server http://docker.for.mac.localhost:8080/r/myapp/publish: (404 Not Found)" source="[server-main.go:206:serverMain()]"
cd services
if [[ "$1" == "local" ]]; then
  echo "Deploying local only"
  fn --verbose deploy --all --app vista --local
else
  fn --verbose deploy --all --app vista
fi
cd ..

# Get rid of any existing minio
docker rm -f minio1 || true

sed  -e "s/APP/$APP/" -e "s/DOCKER_LOCALHOST/$DOCKER_LOCALHOST/" < $PWD/scripts/minio_config.json.tmpl > $PWD/scripts/minio_config.json

sed  -e "s/APP/$APP/" -e "s/DOCKER_LOCALHOST/$DOCKER_LOCALHOST/" < $PWD/scripts/mc.json.tmpl > $PWD/scripts/minio_config.json


# if we want to save data outside of container, add this into line below: -v /tmp/export/minio1:/export
docker run -d -p 9000:9000  --rm --name minio1 \
    -e "MINIO_ACCESS_KEY=$STORAGE_ACCESS_KEY" \
    -e "MINIO_SECRET_KEY=$STORAGE_SECRET_KEY" \
    -v $PWD/scripts/minio_config.json:/root/.minio/config.json minio/minio  server /export
sleep 5

docker run --rm -v $PWD:/mc -w /mc --entrypoint=/bin/sh  -e DEMOACCESSKEY=$STORAGE_ACCESS_KEY -e DEMOSECRETKEY=$STORAGE_SECRET_KEY -e DOCKER_LOCALHOST -e VISTA_MODE minio/mc scripts/setup_minio.sh

pushd scripts
  ./setenv.sh

  if [ $VISTA_MODE = flow ]; then
      ./setup_flow.sh
  fi
popd

docker run --rm -v $PWD:/tmp -w /tmp -e PUBNUB_SUBSCRIBE_KEY treeder/temple public/vista.erb public/vista.html
