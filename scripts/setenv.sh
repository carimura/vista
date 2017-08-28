fn apps config set myapp PUBNUB_PUBLISH_KEY $PUBNUB_PUBLISH_KEY
fn apps config set myapp PUBNUB_SUBSCRIBE_KEY $PUBNUB_SUBSCRIBE_KEY
fn apps config set myapp FUNC_SERVER_URL ${FUNC_SERVER_URL}/r/myapp
fn apps config set myapp MINIO_SERVER_URL $MINIO_SERVER_URL
fn apps config set myapp ACCESS $MINIO_ACCESS_KEY
fn apps config set myapp SECRET $MINIO_ACCESS_KEY
fn apps config set myapp BUCKET oracle-vista-out

cd ../services/alert
fn routes config set myapp /alert TWITTER_CONF_KEY $TWITTER_CONF_KEY
fn routes config set myapp /alert TWITTER_CONF_SECRET $TWITTER_CONF_SECRET
fn routes config set myapp /alert TWITTER_TOKEN_KEY $TWITTER_TOKEN_KEY
fn routes config set myapp /alert TWITTER_TOKEN_SECRET $TWITTER_TOKEN_SECRET

cd ../scraper
fn routes config set myapp /scraper FLICKR_API_KEY $FLICKR_API_KEY
fn routes config set myapp /scraper FLICKR_API_SECRET $FLICKR_API_SECRET
