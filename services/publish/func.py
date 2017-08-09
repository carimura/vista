import json
import sys, os
from pubnub.pnconfiguration import PNConfiguration
from pubnub.pubnub import PubNub

pnconfig = PNConfiguration()
pnconfig.subscribe_key = os.environ["PUBNUB_SUBSCRIBE_KEY"]
pnconfig.publish_key = os.environ["PUBNUB_PUBLISH_KEY"]
pnconfig.ssl = False
pn = PubNub(pnconfig)

def getPayload():
    std_in = sys.stdin.read()
    return json.loads(std_in)

def callback(message):
     print message

def main():
    payload = getPayload()

    bucket_name = payload["Records"][0]["s3"]["bucket"]["name"]
    image_key = payload["Records"][0]["s3"]["object"]["key"]
    url = os.environ["MINIO_SERVER_URL"] + "/" + bucket_name + "/" + image_key

    message = {'url': url, 'id': image_key}
    message_json = json.dumps(message)

    pn.publish().channel(bucket_name).message([message_json]).use_post(True).sync()

if __name__ == "__main__":
    main()
