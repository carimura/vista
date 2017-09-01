set -ex


# we need to have a publish function in place before starting minio, so deploy first
# time="2017-08-30T16:21:13Z" level=error msg="Initializing object layer failed" cause="Unable to initialize event notification. Unexpected response from webhook server http://docker.for.mac.localhost:8080/r/myapp/publish: (404 Not Found)" source="[server-main.go:206:serverMain()]"
cd services
if [[ "$1" == "local" ]]; then
  echo "Deploying local only"
  make deploy-local
else
  make deploy
fi
cd ..

# Get rid of any existing minio
docker rm -f minio1 || true
# if we want to save data outside of container, add this into line below: -v /tmp/export/minio1:/export
docker run -d -p 9000:9000 --name minio1 -v $PWD/scripts/minio_config.json:/root/.minio/config.json minio/minio server /export
sleep 5

docker run --rm -v $PWD:/mc -w /mc --entrypoint=/bin/sh minio/mc scripts/setup_minio.sh

cd scripts
./setenv.sh
cd ..

docker run --rm -v $PWD:/tmp -w /tmp -e PUBNUB_SUBSCRIBE_KEY treeder/temple public/vista.erb public/vista.html
