FROM jfloff/alpine-python:2.7-slim

RUN apk add --update alpine-sdk ca-certificates && rm -rf /var/cache/apk/*
RUN apk add --update py-pillow && rm -rf /var/cache/apk/*
RUN pip install boto3

WORKDIR /func
ADD . /func/
ENTRYPOINT ["python", "func.py"]
