import fdk
import ujson
import os

from fdk import fixtures

from pubnub.pnconfiguration import PNConfiguration
from pubnub.pubnub import PubNub


pnconfig = PNConfiguration()
pnconfig.subscribe_key = os.environ.get("PUBNUB_SUBSCRIBE_KEY")
pnconfig.publish_key = os.environ.get("PUBNUB_PUBLISH_KEY")
pnconfig.ssl = False
pn = PubNub(pnconfig)


TEST_MODE = os.environ.get("TEST_MODE", "false")


# very basic test, will fail if env var TEST_MODE is not set to true
# we don't really want to test pubnub client in unit tests.
async def test_override_content_type(aiohttp_client):
    with open("payload.json", "r") as payload_file:
        call = await fixtures.setup_fn_call(
            aiohttp_client, handle, json=ujson.load(payload_file))
        content, status, headers = await call

        assert 200 == status


def handle(ctx, data=None, **kwargs):
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
                os.environ.get("MINIO_SERVER_URL"),
                bucket_name, image_key)

            message = {'url': url, 'id': image_key}
            message_json = ujson.dumps(message)

            if TEST_MODE not in ['true', '1', 't', 'y', 'yes',
                             'yeah', 'yup', 'certainly', 'uh-huh']:

                (pn.publish().channel(bucket_name).
                 message([message_json]).use_post(True).sync())


if __name__ == "__main__":
    fdk.handle(handle)
