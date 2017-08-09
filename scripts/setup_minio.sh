mc mb myminio/oracle-vista-out
mc policy public myminio/oracle-vista-out
mc events add myminio/oracle-vista-out arn:minio:sqs:us-east-1:1:webhook --events put
