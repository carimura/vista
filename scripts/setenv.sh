APP=${APP:-myapp}

fn apps config set ${APP} PUBNUB_PUBLISH_KEY $PUBNUB_PUBLISH_KEY
fn apps config set ${APP} PUBNUB_SUBSCRIBE_KEY $PUBNUB_SUBSCRIBE_KEY
fn apps config set ${APP} FUNC_SERVER_URL ${FUNC_SERVER_URL}/r/${APP}
fn apps config set ${APP} MINIO_SERVER_URL $MINIO_SERVER_URL
fn apps config set ${APP} COMPLETER_BASE_URL http://$DOCKER_LOCALHOST:8081
fn apps config set ${APP} STORAGE_ACCESS_KEY $STORAGE_ACCESS_KEY
fn apps config set ${APP} STORAGE_SECRET_KEY $STORAGE_SECRET_KEY
fn apps config set ${APP} STORAGE_BUCKET oracle-vista-out
fn apps config set ${APP} FN_TOKEN $FN_TOKEN


cd ../services/alert
fn routes config set ${APP} /alert TWITTER_CONF_KEY $TWITTER_CONF_KEY
fn routes config set ${APP} /alert TWITTER_CONF_SECRET $TWITTER_CONF_SECRET
fn routes config set ${APP} /alert TWITTER_TOKEN_KEY $TWITTER_TOKEN_KEY
fn routes config set ${APP} /alert TWITTER_TOKEN_SECRET $TWITTER_TOKEN_SECRET

cd ../scraper
fn routes config set ${APP} /scraper FLICKR_API_KEY $FLICKR_API_KEY
fn routes config set ${APP} /scraper FLICKR_API_SECRET $FLICKR_API_SECRET

cd ../post-slack
fn routes config set ${APP} /post-slack SLACK_API_TOKEN $SLACK_API_TOKEN

sync_async_fns="alert detect-faces detect-plates draw"

# the flow version requires some functions to be sync
# the normal version requires them to be async
#
if [[ ${VISTA_MODE} == "flow" ]]
then
   echo "-------- Configuring App for Fn Flow ---------"
   # just the flow  version
   fn apps config set ${APP} NO_CHAIN true

   for func in $sync_async_fns ; do
     cd ../$func
     fn routes update  ${APP} $func --type sync
   done
else
   echo "------- Configuring App for Async --------"
   fn apps config set ${APP} NO_CHAIN ""

  for func in $sync_async_fns ; do
    cd ../$func
    fn routes update  ${APP} ${func} --type async
  done
fi
