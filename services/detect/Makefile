deps:
	# need to figure out how to build from something other than a previous version....?
	docker run --rm -v "$(PWD)":/worker -w /worker carimura/detect sh deps.sh

test:
	docker run --rm -e "PAYLOAD_FILE=payload.json" -e "PYTHONPATH=/worker/packages" -v "$(PWD)":/worker -w /worker carimura/detect:$(FACES_VERSION) python /worker/detect.py
