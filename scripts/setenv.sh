#!/usr/bin/env bash

set -ex

APP=${APP:-vista}

fn config app ${APP} PUBNUB_PUBLISH_KEY $PUBNUB_PUBLISH_KEY
fn config app ${APP} PUBNUB_SUBSCRIBE_KEY $PUBNUB_SUBSCRIBE_KEY
fn config app ${APP} FUNC_SERVER_URL ${FUNC_SERVER_URL}/t/${APP}
fn config app ${APP} MINIO_SERVER_URL $MINIO_SERVER_URL
fn config app ${APP} COMPLETER_BASE_URL http://${DOCKER_LOCALHOST}:8081
fn config app ${APP} STORAGE_ACCESS_KEY $STORAGE_ACCESS_KEY
fn config app ${APP} STORAGE_SECRET_KEY $STORAGE_SECRET_KEY
fn config app ${APP} STORAGE_BUCKET oracle-vista-out
fn config app ${APP} STORAGE_BUCKET ${STORAGE_BUCKET:-oracle-vista-out}
fn config app ${APP} S3_REGION ${S3_REGION:-us-phoenix-1}

fn config fn ${APP} alert TWITTER_CONF_KEY $TWITTER_CONF_KEY
fn config fn ${APP} alert TWITTER_CONF_SECRET $TWITTER_CONF_SECRET
fn config fn ${APP} alert TWITTER_TOKEN_KEY $TWITTER_TOKEN_KEY
fn config fn ${APP} alert TWITTER_TOKEN_SECRET $TWITTER_TOKEN_SECRET

fn config f ${APP} scraper-py FLICKR_API_KEY $FLICKR_API_KEY
fn config fn ${APP} scraper-py FLICKR_API_SECRET $FLICKR_API_SECRET

fn config fn ${APP} post-slack SLACK_API_TOKEN $SLACK_API_TOKEN

fn config fn ${APP} flow POST_SLACK_FUNC_ID $(fn inspect fn ${APP} post-slack id | xargs)
fn config fn ${APP} flow SCRAPER_FUNC_ID $(fn inspect fn ${APP} scraper-py id | xargs)
fn config fn ${APP} flow DETECT_PLATES_FUNC_ID $(fn inspect fn ${APP} detect-plates id | xargs)
fn config fn ${APP} flow ALERT_FUNC_ID $(fn inspect fn ${APP} alert id | xargs)
fn config fn ${APP} flow DRAW_FUNC_ID $(fn inspect fn ${APP} draw id | xargs)

############################################################
# async doesn't work for now,
############################################################
fn config app ${APP} NO_CHAIN true
