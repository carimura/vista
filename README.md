# Vista

![logo](https://drive.google.com/uc?export=view&id=0BzyYzhn6bht-Sm1VaFdKY2hKYXc)

## Running Locally

This demo is designed to be run outside of a corporate firewall. The
moving parts have not been instrumented to run with proxies.

### Step 1: Get prerequisite accounts

- Docker ID [docker id docs](https://docs.docker.com/docker-id/)
- Pubnub free account
- Twitter [developer account](https://apps.twitter.com/)
  - NOTE: MAKE A NEW ACCOUNT! Otherwise this will tweet a bunch of things to your Twitter.
- Flickr [developer account](https://www.flickr.com/services/apps/create/apply/)

### Step 2: Set required env vars

1. Copy `envs-example.sh` to `envs.sh` and fill in all the required values.
1. Run: `. ./envs.sh` - NOTE: You NEED the initial `.`

### Step 2: Install Fn CLI and start Fn server

Ensure Docker is running. Ensure you are logged in with docker login.

Install `fn` CLI:

```bash
curl -LSs https://raw.githubusercontent.com/fnproject/cli/master/install | sh
```

Then start `fn` server:

```bash
fn start
```

(Easy huh?)

### Step 3: Set everything up

Use the `local` arg to do everything locally (ie: doesn't push to docker registry).

```bash
./setup.sh [local]
```

### Step 4: Run the demo!

When running in non-flow mode, you can do

```bash
./run.sh
```
This will open a browser window to view the results. 

When running in flow mode, open [the flow ui](http://localhost:3000/#/),
cd to `services/flow` and do

You should also see activity in the server logs, and output to the
vista.html screen. As the draw function finishes, the final images will
push to the screen. Plate detection will also Tweet out from the alert
function.

```bash
cat payload.json | fn call vista flow
```

Make sure your slack link has a channel called `demostream`.

## Future Work Ideas

- All counts will be moved to Functions UI so can simplify by removing those
  graphs.
- Video stream (split-video) added. Could do more work around this. TBD.
- Also see [issues](https://github.com/carimura/vista/issues)

Other ideas: Feel free to create GitHub issues or contact Chad

## Understanding How This Works

In this sort of mashup, the most important thing to understand are the
boxes and the lines. The boxes are the services and the lines are the
wiring between the services. Let's look at the wiring first.

### Wiring

The demo uses minio as a stand-in for Amazon S3 storage. This enables
the demo to be written for the S3 API, yet not have to actually use S3
for the storage.

The functions themselves are set up and "routed" when the `make deploy`
happens.

The act of running the `setenv.sh` script pushes the current set of env
vars into the function runtimes. This includes mashup keys, tokens, and
secrets vital to the success of the demo. 

After this completes, you can do

```bash
$ fn list triggers vista
    FUNCTION        NAME            TYPE    SOURCE          ENDPOINT
    alert           alert           http    /alert          http://localhost:8080/t/vista/alert
    detect-plates   detect-plates   http    /detect-plates  http://localhost:8080/t/vista/detect-plates
    draw            draw            http    /draw           http://localhost:8080/t/vista/draw
    flow            flow            http    /flow           http://localhost:8080/t/vista/flow
    post-slack      post-slack      http    /post-slack     http://localhost:8080/t/vista/post-slack
    publish         publish         http    /publish        http://localhost:8080/t/vista/publish
    scraper         scraper         http    /scraper        http://localhost:8080/t/vista/scraper
```

Finally, the Minio wiring. I'm a bit foggy on this, but it basically
sets up something that looks like an Amazon S3 bucket. 

### Services

#### scraper

This Ruby based service service uses the Flickr API, as exposed to Ruby
by
<[https://github.com/hanklords/flickraw](https://github.com/hanklords/flickraw)>.
The service is built from the Dockerfile, which includes a Gemfile that
pulls in the required dependencies flickraw, json, and rest-client. The
Dockerfile lists func.rb as the entrypoint.

Looking at func.rb, it pulls in the "payload" from stdin, expecting it
to be JSON that looks like

```json
    {
      "query": "license plate car usa",
      "num": "20",
      "countrycode": "us",
      "service_to_call": "detect-plates"
    }
```
The function does the secret sauce call here:

```
    photos = flickr.photos.search(
        :text => search_text,
        :per_page => num_results,
      :page => page,
        :extras => 'original_format',
        :safe_search => 1,
        :content_type => 1
    )
```

For each photo in the result set, do a POST to
FUNC_SERVER_URL/detect-plates. Post data is the following:

```
    payload = {:id => photo.id, 
               :image_url => image_url,
               :countrycode => payload_in["countrycode"],
               :bucket => payload_in["bucket"]
    }
```


#### detect-plates

This service is a go based service that uses two dependencies, as well
as a whole passel of standard go packages.

* The OpenALPR (Automatic License Plate Recognition), with the go
language binding.
[https://github.com/openalpr/openalpr](https://github.com/openalpr/openalpr)

* The pubnub messaging service, with the go language binding [https://github.com/pubnub/go](https://github.com/pubnub/go)

Looking at the main.go, it looks like this service reads the POST
payload, calls fnstart(), passing the BUCKET env var and the photo id.
Simply appears to fob the request off to pubnub. I think this may be a
way to allow other services to asynchronously consume the result of this
service's plate detection. The secret sauce function is here:

```go
results, err := alpr.RecognizeByBlob(imageBytes)
```

For each entry in results, build another payload, adding a "rectangles"
array and the text of the plate number. The payload looks like:

```go
    pout := &payloadOut{
        ID:         p.ID,
        ImageURL:   p.URL,
        Rectangles: []rectangle{{StartX: plate.PlatePoints[0].X, StartY: plate.PlatePoints[0].Y, EndX: plate.PlatePoints[2].X, EndY: plate.PlatePoints[2].Y}},
        Plate:      plate.BestPlate,
    }
```

This is posted further downstream to the draw function at

```go
postURL := os.Getenv("FUNC_SERVER_URL") + "/draw"
```

and also a copy of the same payload is posted to the alert service.

```go
alertPostURL := os.Getenv("FUNC_SERVER_URL") + "/alert"
```

#### Draw service

We're back in Ruby for this one. But the runtime of the func.yaml is ""
instead of ruby. This means to use the Dockerfile to build the
function. This ruby service uses alpine-sdk and imagemagick. Upon
receipt of a payload, first we send a message to pubnub saying we are
running an image. The image is downloaded from the payload. We use
ImageMagick to draw the rectangles in the payload on the image. We then
upload the new image to MINIO, which looks like it is acting as a
standin to Amazon S3. Finally we send a message to pubnub saying this
run is done.

#### Alert service

This one is in Go again. Runtime is the empty string. Again, is this
picked up from the func.go?  This uses the Anaconda go twitter client
library
[https://github.com/ChimeraCoder/anaconda](https://github.com/ChimeraCoder/anaconda). It also uses pubnub.

Recall that the payload coming to us is the same one sent to the image
service. It downloads the image from the URL, converts it to base64,
and composes a tweet including the image and the plate text. As with
detect-plates, this uses pubnub to update the status.

#### vista.html

It looks like this entirely relies on pubnub to draw the fancy
dynamically updating line chart and the images as they come along.

This uses the pubnub javascript API, as well as highcharts.




