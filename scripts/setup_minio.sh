mc mb local/oracle-vista-out
mc policy public local/oracle-vista-out
mc events add local/oracle-vista-out arn:minio:sqs:us-east-1:1:webhook --events put
