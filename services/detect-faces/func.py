import sys, os
import cv
import json
import urllib2
import Algorithmia
import requests

std_in = sys.stdin.read()
payload = json.loads(std_in)

no_chain = False
if os.environ["NO_CHAIN"]:
    no_chain = True

def detectObjects(image):
    grayscale = cv.CreateImage(cv.GetSize(image), 8, 1)
    cv.CvtColor(image, grayscale, cv.CV_BGR2GRAY)

    cv.EqualizeHist(grayscale, grayscale)
    cascade = cv.Load(os.getcwd() + '/haarcascade_frontalface_alt.xml')
    # Haar feature-based cascade classifiers
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


def downloadFile(url):
    u = urllib2.urlopen(url)
    filename = payload["id"] + '.jpg'

    localFile = open(filename, 'w')
    localFile.write(u.read())
    localFile.close()
    return filename

# I'm told to use this one instead:
# https://algorithmia.com/algorithms/sfw/NudityDetectioni2v
def isNude(url):
    test = "http://www.isitnude.com.s3-website-us-east-1.amazonaws.com/assets/images/sample/young-man-by-the-sea.jpg"
    client = Algorithmia.client(payload["algorithmia_key"])
    algo = client.algo('sfw/NudityDetectioni2v/0.2.12')

    #url = test
    print >> sys.stderr, "is_nude url: " + url
    
    response = algo.pipe(str(url))
    print >> sys.stderr, "is_nude response: " + str(response)
    
    result = response.result
    print >> sys.stderr, "is_nude result: " + str(result)
    print >> sys.stderr, "is_nude true/false: " + str(result["nude"])
    
    return result["nude"]
    

def main():
    # Notify the UI that a function has started
    image_name = payload["id"] + ".jpg"

    image_url = payload["image_url"]
    print >> sys.stderr, "image_url: " + image_url

    f = downloadFile(image_url)
    image = cv.LoadImageM(image_name)

    is_nude = False #isNude(image_url)
    
    cat_url = "http://random.cat/meow"
    
    if is_nude:
       cat_req = requests.get(cat_url)
       cat_json = cat_req.json()
       print >> sys.stderr, "cat_json: " + str(cat_json)
       image_url = cat_json["file"]

    rectangles = detectObjects(image)

    print >> sys.stderr, "rectangles: " + str(rectangles)

    next_payload = {
      "image_url": image_url,
      "is_nude": is_nude,
      "rectangles": rectangles,
      "id": payload["id"]
    }

    if no_chain:
       print json.dumps(next_payload)
    else:
       post_url = os.environ["FUNC_SERVER_URL"] + "/draw"
       r = requests.post(post_url, data=json.dumps(next_payload))


if __name__ == "__main__":
    main()
