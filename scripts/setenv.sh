#!/usr/bin/env bash

set -ex

APP=${APP:-vista}

fn config app ${APP} PUBNUB_PUBLISH_KEY $PUBNUB_PUBLISH_KEY
fn config app ${APP} PUBNUB_SUBSCRIBE_KEY $PUBNUB_SUBSCRIBE_KEY
fn config app ${APP} FUNC_SERVER_URL ${FUNC_SERVER_URL}/r/${APP}
fn config app ${APP} MINIO_SERVER_URL $MINIO_SERVER_URL
fn config app ${APP} COMPLETER_BASE_URL $COMPLETER_BASE_URL
fn config app ${APP} STORAGE_ACCESS_KEY $STORAGE_ACCESS_KEY
fn config app ${APP} STORAGE_SECRET_KEY $STORAGE_SECRET_KEY
fn config app ${APP} STORAGE_BUCKET oracle-vista-out
fn config app ${APP} FN_TOKEN $FN_TOKEN


cd ../services/alert
fn config fn ${APP} alert TWITTER_CONF_KEY $TWITTER_CONF_KEY
fn config fn ${APP} alert TWITTER_CONF_SECRET $TWITTER_CONF_SECRET
fn config fn ${APP} alert TWITTER_TOKEN_KEY $TWITTER_TOKEN_KEY
fn config fn ${APP} alert TWITTER_TOKEN_SECRET $TWITTER_TOKEN_SECRET

cd ../scraper
fn config f ${APP} scraper FLICKR_API_KEY $FLICKR_API_KEY
fn config fn ${APP} scraper FLICKR_API_SECRET $FLICKR_API_SECRET

cd ../post-slack
fn config fn ${APP} post-slack SLACK_API_TOKEN $SLACK_API_TOKEN



############################################################
# async doesn't work for now,
############################################################
# the flow version requires some functions to be sync
# the normal version requires them to be async
#
# sync_async_fns="alert detect-plates draw"
#if [[ ${VISTA_MODE} == "flow" ]]
#then
#   echo "-------- Configuring App for Fn Flow ---------"
#   # just the flow  version
#   fn config app ${APP} NO_CHAIN true
#
#   for func in $sync_async_fns ; do
#     cd ../$func
#     fn update route ${APP} $func --type sync
#   done
#else
#   echo "------- Configuring App for Async --------"
#   fn config app ${APP} NO_CHAIN ""
#
#  for func in $sync_async_fns ; do
#    cd ../$func
#    fn update route ${APP} ${func} --type async
#  done
#fi
