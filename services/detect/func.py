import sys, os
import cv
import json
import urllib2
import Algorithmia
import requests

from pubnub.pnconfiguration import PNConfiguration
from pubnub.pubnub import PubNub

def detectObjects(image):
    grayscale = cv.CreateImage(cv.GetSize(image), 8, 1)
    cv.CvtColor(image, grayscale, cv.CV_BGR2GRAY)

    cv.EqualizeHist(grayscale, grayscale)
    cascade = cv.Load(os.getcwd() + '/haarcascade_frontalface_alt.xml')
    faces = cv.HaarDetectObjects(grayscale, cascade, cv.CreateMemStorage(), 1.2, 2, cv.CV_HAAR_DO_CANNY_PRUNING, (20,20))
    rectangles = []
    if faces:
        for f in faces:
            rectangles.append({
              "startx": f[0][0],
              "starty": f[0][1],
              "endx": f[0][0] + f[0][2],
              "endy": f[0][1] + f[0][3]
            })
    return rectangles


def getPayload():
  return json.loads(sys.stdin.read())

def downloadFile(url, payload_id):
    u = urllib2.urlopen(url)
    filename = payload_id + '.jpg'

    localFile = open(filename, 'w')
    localFile.write(u.read())
    localFile.close()
    return filename

def sendWorkerCount(bucket_name, image_key):
    print "sendWorkerCount image_key: " + str(image_key)    
    fixed_image_key = "image_" + image_key
    message = {'id': fixed_image_key}
    message_json = json.dumps(message)
    print "sendWorkerCount message_json: " + str(message_json)
    pn = PubNub(pnconfig)
    pn.publish().channel(bucket_name).message([message_json]).use_post(True).sync()

# I'm told to use this one instead:
# https://algorithmia.com/algorithms/sfw/NudityDetectionEnsemble
def isNude(url):
    test = "http://www.isitnude.com.s3-website-us-east-1.amazonaws.com/assets/images/sample/young-man-by-the-sea.jpg"
    client = Algorithmia.client('sim0b6hAy8ZVhgM4MAh6xfcbcjo1')
    algo = client.algo('sfw/NudityDetection/1.1.0')

    #url = test
    print "is_nude url: " + url
    
    response = algo.pipe(str(url))
    print "is_nude response: " + str(response)
    
    result = response.result
    print "is_nude result: " + str(result)
    print "is_nude true/false: " + result["nude"]
    
    return result["nude"]
    

def main():
    payload = getPayload()
    print "PAYLOAD: " + str(payload)

    pnconfig = PNConfiguration()
    pnconfig.subscribe_key = payload["pubnub_subscribe_key"]
    pnconfig.publish_key = payload["pubnub_publish_key"]
    pnconfig.ssl = False
    pubnub = PubNub(pnconfig)

    # Notify the UI that a function has started
    image_name = payload["id"] + ".jpg"
    sendWorkerCount("oracle-faces-out", image_name)

    if "Message" in payload:
      print "feels like SNS"
      message_json = json.loads(payload["Message"])
      bucket_name = message_json["Records"][0]["s3"]["bucket"]["name"]
      print "Bucket: " + bucket_name
      image_key = message_json["Records"][0]["s3"]["object"]["key"]
      print "image_key: " + image_key
      image_url = "https://s3.amazonaws.com/"+bucket_name+"/" + image_key
    else:
      print "seems like direct queue from image url"
      image_url = payload["image_url"]

    print "image_url: " + image_url

    f = downloadFile(image_url, payload["id"])
    image = cv.LoadImageM(image_name)

    is_nude = "false"#isNude(image_url)

    rectangles = detectObjects(image)

    print "rectangles: " + str(rectangles)

    next_payload = {
      "image_url": image_url,
      "is_nude": is_nude,
      "rectangles": rectangles,
      "aws_access": payload["aws_access"],
      "aws_secret": payload["aws_secret"],
      "aws_s3_bucket_name": payload["aws_bucket"],
      "id": payload["id"]
    }
    
    post_url = payload["func_server_url"] + "/draw"

    r = requests.post(post_url, data=json.dumps(next_payload))


if __name__ == "__main__":
    main()
