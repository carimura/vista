FROM python:2.7.13-alpine

RUN apk add --update python-dev
RUN apk add --update py-pip
RUN apk add --update alpine-sdk

RUN pip install 'pubnub>=4.0.12'

WORKDIR /func
ADD . /func/

ENTRYPOINT ["python", "func.py"]

