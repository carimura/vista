set -ex

docker run --rm -v $PWD:/mc -w /mc --entrypoint=/bin/sh minio/mc scripts/setup_minio.sh

cd services
make deploy
cd ..

cd scripts
./setenv.sh
cd ..

cd public
docker run --rm -v $PWD:/tmp -w /tmp -e PUBNUB_SUBSCRIBE_KEY treeder/temple vista.erb vista.html
cd ..
