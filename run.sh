#!/usr/bin/env bash

set -ex

open public/vista.html

cd services/scraper
cat payload.json | fn invoke ${APP:-vista} scraper
cd ../..
