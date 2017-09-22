mc mb fnstage/oracle-vista-out
mc mb fnstage/videoimages
mc policy public fnstage/oracle-vista-out
mc policy public fnstage/videoimages

mc events add fnstage/oracle-vista-out arn:minio:sqs:us-east-1:1:webhook --events put || true
