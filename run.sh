set -ex

open public/vista.html

cd services/scraper
cat payload.json | fn call ${APP:-vista} /scraper
cd ../..
