#!/usr/bin/env bash

set -ex

mc config host add local http://docker.for.mac.localhost:9000  $DEMOACCESSKEY $DEMOSECRETKEY

mc mb local/oracle-vista-out
mc mb local/videoimages
mc policy public local/oracle-vista-out
mc policy public local/videoimages

mc events add local/oracle-vista-out arn:minio:sqs:us-east-1:1:webhook --events put || true
