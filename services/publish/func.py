import fdk
import ujson
import os

from pubnub.pnconfiguration import PNConfiguration
from pubnub.pubnub import PubNub


pnconfig = PNConfiguration()
pnconfig.subscribe_key = os.environ.get("PUBNUB_SUBSCRIBE_KEY")
pnconfig.publish_key = os.environ.get("PUBNUB_PUBLISH_KEY")
pnconfig.ssl = False
pn = PubNub(pnconfig)


def main(ctx, data=None, **kwargs):
    payload = None
    if data and len(data) > 0:
        payload = ujson.loads(data)

    records = payload.get("Records", [])

    for record in records:
        s3 = record.get("s3", {})
        obj, bucket = s3.get("object", {}), s3.get("bucket", {})
        bucket_name = bucket.get("name")
        image_key = obj.get("key")
        if all((bucket_name, image_key)):
            url = "{0}/{1}/{2}".format(
                os.environ["MINIO_SERVER_URL"],
                bucket_name, image_key)

            message = {'url': url, 'id': image_key}
            message_json = ujson.dumps(message)

            (pn.publish().channel(bucket_name).
             message([message_json]).use_post(True).sync())


if __name__ == "__main__":
    fdk.handle(main)
