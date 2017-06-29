Oracle Faces

Boom demo.....


1.  make sure $API_URL is set
2. `cd services`
3. `make deploy` (this should deploy detect, draw, publish to fn server)
4. `cd scraper`
5. `mv payload.sample.json payload.json` 
6. enter various creds into payload.json
7. `cat payload.json | fn run`

there will be no UI yet until I explain how to setup pubnub channel....

but we can at least get started here.
