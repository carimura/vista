Oracle Vista

1.  make sure $API_URL is set
1. `cd services`
1. `make deploy` (this should deploy both detects, draw, publish, and alert to FN server)
1. cd ..; ./setenv.sh (this will set all the proper function secrets, get from
   Chad)
1. `cd public; open vista.html`
1. Enter "oracle-vista-out" into the box (this subscribes to pubnub channel)
1. `cd scraper`
1. `cat payload.json | fn call myapp /scraper`

there will be no UI yet until I explain how to setup pubnub channel....

but we can at least get started here.
