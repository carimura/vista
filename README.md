# Vista App

## setup the app

1.  make sure $API_URL is set
1. `cd services`
1. `make deploy` (this should deploy both detects, draw, publish, and alert to FN server)
1. cd ..; 
1. fill out setenv.sh then run ./setenv.sh (this will set all the proper function secrets)
1. `cd public; open vista.html`
1. Enter "oracle-vista-out" into the box (this subscribes to pubnub channel)
1. `cd scraper`
1. `cat payload.json | fn call myapp /scraper`



## minio

The app needs Minio to run somewhere since that is the S3-compliant storage
engine. Also the webhooks need to be configured so that they can push out the
results to the publish function which pushes the images to the
public/vista.html front end.

1. Make sure Minio is setup and the bucket is created.
1. The minio config must have webhooks setup per [this blog post](https://blog.minio.io/introducing-webhooks-for-minio-e2c3ad26deb2)
1. After you create the webhooks, you'll get ARN's, from there you need to add those arns to bucket events like so:
- `mc events add fnstage/oracle-vista-out arn:minio:sqs:us-east-1:2:webhook --events put`
- `mc events add fnstage/oracle-vista-out arn:minio:sqs:us-east-1:2:webhook --events put`
1. That should be it for Minio.


