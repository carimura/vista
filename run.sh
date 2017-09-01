set -ex

open public/vista.html

cd services/scraper
cat payload.json | fn call myapp /scraper
# 1. Facial detection: `cat payload_faces.json | fn call myapp /scraper`
cd ../..
