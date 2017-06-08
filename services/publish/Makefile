deps:
	docker run --rm -v "$(PWD)":/func -w /func funcy/python:2.7.10-dev sh deps.sh

test:
	docker run --rm -e "PAYLOAD_FILE=payload.json" -e "PYTHONPATH=/func/packages" -v "$(PWD)":/func -w /func funcy/python:2.7.10 python /func/publish.py
