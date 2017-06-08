import sys, os
import cv
import json
import urllib2
import Algorithmia

from pubnub.pnconfiguration import PNConfiguration
from pubnub.pubnub import PubNub

pnconfig = PNConfiguration()
pnconfig.subscribe_key = "sub-1e453968-bc05-11e0-9cf9-cbaf6932e4b8"
pnconfig.publish_key = "pub-025536de-c773-415a-9961-3d5c2bec5f26"
pnconfig.ssl = False
 
pubnub = PubNub(pnconfig)

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


def getTaskID():
  return "12345"


def downloadFile(url):
    u = urllib2.urlopen(url)
    localFile = open(getTaskID()+'.jpg', 'w')
    localFile.write(u.read())
    localFile.close()
    return getTaskID()+'.jpg'

def callback(message):
     print message

def sendWorkerCount(bucket_name, image_key):
    print "sendWorkerCount image_key: " + str(image_key)    
    message = {'id': image_key}
    message_json = json.dumps(message)
    print "sendWorkerCount message_json: " + str(message_json)

    # push data onto pubnub channel
    pn = PubNub(pnconfig)
    pn.publish().channel(bucket_name).message([message_json]).use_post(True).async(callback)

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

    # Notify the UI that a worker has started
    image_id = getTaskID()+".jpg"
    sendWorkerCount("iron-faces-out", image_id)

    if "Message" in payload:
      message_json = json.loads(payload["Message"])
      print "seems like SNS"
      bucket_name = message_json["Records"][0]["s3"]["bucket"]["name"]
      print "Bucket: " + bucket_name
      image_key = message_json["Records"][0]["s3"]["object"]["key"]
      print "image_key: " + image_key
      image_url = "https://s3.amazonaws.com/"+bucket_name+"/" + image_key
    else:
      print "seems like direct queue from image url"
      image_url = payload["image_url"]

    print "image_url: " + image_url

    f = downloadFile(image_url)
    image = cv.LoadImageM(image_id)
    #is_nude = "false"#isNude(image_url)
    rectangles = detectObjects(image)

    print "rectangles: " + str(rectangles)

    print str(os.environ)

    #next_payload = {
    #  "image_url": image_url,
    #  "is_nude": is_nude,
    #  "rectangles": rectangles,
    #  "aws_access": "AKIAJFFRXQSLDGC3EGAQ",
    #  "aws_secret": "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",
    #  "aws_s3_bucket_name": "iron-faces-out",
    #  "id": getTaskID()
    #}
    #worker = IronWorker()
    #task = Task(code_name="carimura/draw", payload=next_payload)
    #worker.queue(task)

if __name__ == "__main__":
    main()
