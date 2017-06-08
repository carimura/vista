
clean:
	cd detect;rm -Rf packages;
	cd draw; rm -Rf bundle .bundle;
	cd publish; rm -Rf packages;
	cd scraper; rm -Rf bundle .bundle;

build:
	cd detect;make deps; docker build -t carimura/detect:$(FACES_VERSION) .
	cd draw;make deps; docker build -t carimura/draw:$(FACES_VERSION) .
	cd publish;make deps; docker build -t carimura/publish:$(FACES_VERSION) .
	cd scraper;make deps; docker build -t carimura/scraper:$(FACES_VERSION) .

push:
	docker push carimura/scraper:$(FACES_VERSION)
	docker push carimura/detect:$(FACES_VERSION)
	docker push carimura/draw:$(FACES_VERSION)
	docker push carimura/publish:$(FACES_VERSION)

register:

run:
