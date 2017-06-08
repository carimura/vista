import json
import pubnub
from iron_worker import *

def getPayload():
    with open(os.environ['PAYLOAD_FILE'], 'r') as f:
        data = f.read()
        return json.loads(data)

def callback(message):
     print message

def main():
    # get payload
    payload = getPayload()
    message_json = json.loads(payload["Message"])

    # extract data
    bucket_name = message_json["Records"][0]["s3"]["bucket"]["name"]
    print "Bucket: " + bucket_name

    image_key = message_json["Records"][0]["s3"]["object"]["key"]
    print "Image Key: " + image_key

    url = "https://s3.amazonaws.com/"+bucket_name+"/" + image_key;
    
    # build and package message as json
    message = {'url': url, 'id': image_key}
    message_json = json.dumps(message)

    # build pubnub request
    pn = pubnub.Pubnub(publish_key='pub-b6bcf2ef-b6a6-45f4-a974-5640cf2b50f7', subscribe_key='sub-3e0a8b35-9316-11e1-910f-a9e1ad10d598')

    # push data onto pubnub channel
    pn.publish(bucket_name, message_json, callback=callback, error=callback)

if __name__ == "__main__":
    main()
