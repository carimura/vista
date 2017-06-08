#!/bin/sh

apk add --update python-dev
apk add --update py-pip
apk add --update alpine-sdk
pip install -t packages -r requirements.txt
