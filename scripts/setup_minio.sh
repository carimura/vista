mc mb fnstage/oracle-vista-out
mc policy public fnstage/oracle-vista-out
mc events add fnstage/oracle-vista-out arn:minio:sqs:us-east-1:1:webhook --events put
mc events add fnstage/oracle-vista-out arn:minio:sqs:us-east-1:2:webhook --events put

