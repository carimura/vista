# Vista App

![logo](https://drive.google.com/uc?export=view&id=0BzyYzhn6bht-Sm1VaFdKY2hKYXc)

## Running Locally

This demo is designed to be run outside of a corporate firewall.  The
moving parts have not been instrumented to run with proxies.

### Step 1: Get prerequisite accounts

- Docker ID [docker id docs](https://docs.docker.com/docker-id/)
- Pubnub free account
- Twitter [developer account](https://apps.twitter.com/)
- Flickr [developer account](https://www.flickr.com/services/apps/create/apply/)

### Step 2: Install Fn CLI and start Fn server

Ensure Docker is running.  Ensure you are logged in with docker login.

`curl -LSs https://raw.githubusercontent.com/fnproject/cli/master/install | sh`

Installs to /usr/local/bin/fn.

`fn start`

(Easy huh?)

### Step 3: Setup ngrok
1. install [ngrok](https://ngrok.com/)
1. `ngrok http 8080` (for Fn)
1. `ngrok http 9000` (for minio)
1. `export API_URL=<ngrok_url_for_http_8080>`


### Step 4: Deploy/configure the Vista functions
Ensure you have a GNU compatible make.

1. modify every func.yaml in /services/\* to change carimura to your docker id
   (working on a better way to manage this, but in the meantime `find . -name func.yaml -exec perl -pi.bak -e "s/carimura/<yourdockerid>/g" {} \; -print`)
1. `cd services; make deploy` (this should deploy all demo funcs to the Fn server) 
1. set the proper ENV vars needed in scripts/setenv.sh, then run run `./setenv.sh`
1. `open public/vista.html`
1. Enter oracle-vista-out as the value of the BUCKET environment variable into the box (this subscribes to pubnub channel).

### Step 5: Local minio setup
1. install the [mc minio client](https://github.com/minio/mc)
1. Edit scripts/minio_config.json to change the webhook URL to the API_URL from step 4 above
1. `mkdir -p /tmp/config/minio1; cp minio_config.json /tmp/config/minio1/config.json`
1. start the minio server 
```
docker run -p 9000:9000 --name minio1 \
   -v /tmp/export/minio1:/export \
   -v /tmp/config/minio1:/root/.minio \
   minio/minio server /export
```
1. `mc config host add local <insert ngrok url for port 9000> DEMOACCESSKEY DEMOSECRETKEY` (local is an alias, can be anything but below script assumes local)
1. Make sure the bucket in scripts/setup_minio.sh is the same as set in step 5, then execute `./setup_minio.sh`. This sets up the bucket and webhooks.

### Step 6: Run the demo!
1. `cd scraper`
1. Plate detection: `cat payload.json | fn call myapp /scraper`
1. Facial detection: `cat payload_faces.json | fn call myapp /scraper`

You should see activity in the ngrok logs, server logs, and output to the vista.html screen. As the draw function finishes, the final images will push to the screen. Plate detection will also Tweet out from the alert function.


## Known Issues

- [Issue 13](https://github.com/carimura/vista/issues/13):  Ngrok gives a 429 Too Many Requests when running this at any scale (ie > 5 images). Upgrading Ngrok fixes this, but we shouldn't force all users of this demo to upgrade ngrok. Maybe future work idea is to get rid of ngrok somehow..

## Future Work Ideas

- Add the function count stuff to the publish function
- Add video stream as source for detect-plates, not just the flickr scraper
- lots of dependencies still. Could eliminate Flickr by pulling straight from Google images?
- if BMC has anything S3-compat, could use that as well to avoid Minio

Other ideas: Feel free to create GitHub issues or contact Chad

## Understanding How This Works

### Services

#### scraper

This Ruby based service service uses the Flickr API, as exposed to Ruby
by
<[https://github.com/hanklords/flickraw](https://github.com/hanklords/flickraw)>.
The service is built from the Dockerfile, which includes a Gemfile that
pulls in the required dependencies flickraw, json, and rest-client.  The
Dockerfile lists func.rb as the entrypoint.

Looking at func.rb, it pulls in the "payload" from stdin, expecting it
to be JSON that looks like

    {
      "query": "license plate car usa",
      "num": "20",
      "countrycode": "us",
      "service_to_call": "detect-plates"
      }

The function does the secret sauce call here:

    photos = flickr.photos.search(
        :text => search_text,
        :per_page => num_results,
      :page => page,
        :extras => 'original_format',
        :safe_search => 1,
        :content_type => 1
    )

For each photo in the result set, do a POST to
FUNC_SERVER_URL/detect-plates.  Post data is the following:

    payload = {:id => photo.id, 
               :image_url => image_url,
               :countrycode => payload_in["countrycode"],
               :bucket => payload_in["bucket"]
    }


#### detect-plates

This service is a go based service that uses two dependencies, as well
as a whole passel of standard go packages.

* The OpenALPR (Automatic License Plate Recognition), with the go
language binding.
[https://github.com/openalpr/openalpr](https://github.com/openalpr/openalpr)

* The pubnub messaging service, with the go language binding [https://github.com/pubnub/go](https://github.com/pubnub/go)

Looking at the main.go, it looks like this service reads the POST
payload, calls fnstart(), passing the BUCKET env var and the photo id.
Simply appears to fob the request off to pubnub.  I think this may be a
way to allow other services to asynchronously consume the result of this
service's plate detection.  The secret sauce function is here:

    results, err := alpr.RecognizeByBlob(imageBytes)

For each entry in results, build another payload, adding a "rectangles"
array and the text of the plate number.  The payload looks like:

    pout := &payloadOut{
        ID:         p.ID,
        ImageURL:   p.URL,
        Rectangles: []rectangle{{StartX: plate.PlatePoints[0].X, StartY: plate.PlatePoints[0].Y, EndX: plate.PlatePoints[2].X, EndY: plate.PlatePoints[2].Y}},
        Plate:      plate.BestPlate,
    }

This is posted further downstream to the draw function at

    postURL := os.Getenv("FUNC_SERVER_URL") + "/draw"

and also a copy of the same payload is posted to the alert service.

    alertPostURL := os.Getenv("FUNC_SERVER_URL") + "/alert"

#### Draw service

We're back in Ruby for this one.  But the runtime of the func.yaml is ""
instead of ruby.  Not sure why that is.  Some defaults thing?  This ruby
service uses alpine-sdk and imagemagick.  Upon receipt of a payload,
first we send a message to pubnub saying we are running an image.  The
image is downloaded from the payload.  We use ImageMagick to draw the
rectangles in the payload on the image.  We then upload the new image to
MINIO, which looks like it is acting as a standin to Amazon S3.  Finally
we send a message to pubnub saying this run is done.

#### Alert service

This one is in go again.  Runtime is the empty string.  Again, is this
picked up from the func.go?  This uses the Anaconda go twitter client
library
[https://github.com/ChimeraCoder/anaconda](https://github.com/ChimeraCoder/anaconda).  It also uses pubnub.

Recall that the payload coming to us is the same one sent to the image
service.  It downloads the image from the URL, converts it to base64,
and composes a tweet including the image and the plate text.  As with
detect-plates, this uses pubnub to update the status.

#### vista.html

It looks like this entirely relies on pubnub to draw the fancy
dynamically updating line chart and the images as they come along.

This uses the pubnub javascript API, as well as highcharts.


