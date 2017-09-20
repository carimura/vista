mc config host add local http://$DOCKER_LOCALHOST:9000 DEMOACCESSKEY DEMOSECRETKEY

mc mb local/oracle-vista-out
mc mb local/videoimages
mc policy public local/oracle-vista-out
mc policy public local/videoimages

# don't put the webapp into the flow config

if [ "${VISTA_MODE}" != "flow" ]; then
  mc events add local/oracle-vista-out arn:minio:sqs:us-east-1:1:webhook --events put || true
fi
