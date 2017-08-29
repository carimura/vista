set -ex

docker run --rm -d -it -p 9000:9000 --name minio1 -v /tmp/export/minio1:/export -v $PWD/scripts/minio_config.json:/root/.minio/config.json minio/minio server /export

docker run --rm -v $PWD:/mc -w /mc --entrypoint=/bin/sh minio/mc scripts/setup_minio.sh

cd services
make deploy
cd ..

cd scripts
./setenv.sh
cd ..

docker run --rm -v $PWD:/tmp -w /tmp -e PUBNUB_SUBSCRIBE_KEY treeder/temple public/vista.erb public/vista.html
