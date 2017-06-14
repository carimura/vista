
#build:
#	cd detect;make deps; docker build -t carimura/detect:$(FACES_VERSION) .
#	cd draw;make deps; docker build -t carimura/draw:$(FACES_VERSION) .
#	cd publish;make deps; docker build -t carimura/publish:$(FACES_VERSION) .
#	cd scraper;make deps; docker build -t carimura/scraper:$(FACES_VERSION) .

deploy:
	cd scraper; fn deploy myapp
	cd draw; fn deploy myapp
	cd publish; fn deploy myapp
	cd detect; fn deploy myapp
