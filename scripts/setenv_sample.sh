fn apps config set myapp PUBNUB_PUBLISH_KEY yourvalue
fn apps config set myapp PUBNUB_SUBSCRIBE_KEY yourvalue
fn apps config set myapp FUNC_SERVER_URL yourvalue
fn apps config set myapp MINIO_SERVER_URL yourvalue
fn apps config set myapp BUCKET oracle-vista-out
fn apps config set myapp ACCESS DEMOACCESSKEY
fn apps config set myapp SECRET DEMOSECRETKEY

cd ../services/alert
fn routes config set myapp /alert TWITTER_CONF_KEY yourvalue
fn routes config set myapp /alert TWITTER_CONF_SECRET yourvalue
fn routes config set myapp /alert TWITTER_TOKEN_KEY yourvalue
fn routes config set myapp /alert TWITTER_TOKEN_SECRET yourvalue

cd ../scraper
fn routes config set myapp /scraper FLICKR_API_KEY yourvalue
fn routes config set myapp /scraper FLICKR_API_SECRET yourvalue


# For testing with fn run. Not needed for running deployed demo per README.md
export PUBNUB_PUBLISH_KEY=yourvalue
export PUBNUB_SUBSCRIBE_KEY=tiyrvalue
export FUNC_SERVER_URL=yourvalue
export MINIO_SERVER_URL=yourvalue
export BUCKET=oracle-vista-out
export ACCESS=DEMOACCESSKEY
export SECRET=DEMOSECRETKEY
export TWITTER_CONF_KEY=yourvalue
export TWITTER_CONF_SECRET=yourvalue
export TWITTER_TOKEN_KEY=yourvalue
export TWITTER_TOKEN_SECRET=yourvalue
export FLICKR_API_KEY=yourvalue
export FLICKR_API_SECRET=yourvalue
