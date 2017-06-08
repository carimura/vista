#!/bin/sh
ln -s /dev/null /dev/raw1394
export PAYLOAD_FILE=/mnt/task/task_payload.json;
python func.py;
