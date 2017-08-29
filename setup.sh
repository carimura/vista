set -ex

docker run --rm -v $PWD:/mc -w /mc --entrypoint=/bin/sh minio/mc scripts/setup_minio.sh

cd services
make deploy
cd ..

cd scripts
./setenv.sh
cd ..
