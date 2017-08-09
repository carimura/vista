# Vista App

## Running Locally

### Step 1: Get prerequisite accounts

- Pubnub free account
- Twitter [developer account](https://apps.twitter.com/)
- Flickr [developer account](https://www.flickr.com/services/apps/create/apply/)

### Step 2: Install Fn CLI and start Fn server
`curl -LSs https://raw.githubusercontent.com/fnproject/cli/master/install | sh`

`fn start`

(Easy huh?)

### Step 3: Setup ngrok
1. install [ngrok](https://ngrok.com/)
1. `ngrok http 8080` (for Fn)
1. `ngrok http 9000` (for minio)
1. `export API_URL=<ngrok_url_for_http_8080>`


### Step 4: Deploy/configure the Vista functions
1. `cd services; make deploy` (this should deploy all demo funcs to the Fn server) 
1. `cp scripts/setenv_sample.sh setenv.sh`, fill out all "yourvalue" values (the other ones should work) then `./setenv.sh`
1. `open public/vista.html`
1. Enter the value you're using for your BUCKET (default: oracle-vista-out) environment variable into the box (this subscribes to pubnub channel)

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
1. `mc config host add myminio <insert ngrok url for port 9000> DEMOACCESSKEY DEMOSECRETKEY`
1. Make sure the bucket in scripts/setup_minio.sh is the same as set in step 5, then execute `./setup_minio.sh`. This sets up the bucket and webhooks.

### Step 6: Run the demo!
1. `cd scraper`
1. `cat payload.json | fn call myapp /scraper`



## Known Issues

- Minio gives a 429 Too Many Requests when firing the webhook under a tiny amount of load. WTF Minio.

## Future Work Ideas

- Add the function count stuff to the publish function
- Add video stream as source for detect-plates, not just the flickr scraper
- lots of dependencies still. Could eliminate Flickr by pulling straight from Google images?

Other ideas: Feel free to create GitHub issues or contact Chad