FROM fnproject/python:3.7-dev as build-stage

ARG PUBNUB_SUBSCRIBE_KEY
ARG PUBNUB_PUBLISH_KEY
ARG MINIO_SERVER_URL

ENV PUBNUB_SUBSCRIBE_KEY=${PUBNUB_SUBSCRIBE_KEY}
ENV PUBNUB_PUBLISH_KEY=${PUBNUB_PUBLISH_KEY}
ENV MINIO_SERVER_URL=${MINIO_SERVER_URL}

WORKDIR /function
ADD . /function/

RUN pip3 install --target /python/  --no-cache --no-cache-dir -r requirements.txt pytest

RUN pip3 install --no-cache --no-cache-dir -r requirements.txt && pytest -v -s --tb=long func.py

FROM fnproject/python:3.7
WORKDIR /function
COPY --from=build-stage /function /function
COPY --from=build-stage /python /python
ENV PYTHONPATH=/python
ENTRYPOINT ["python3", "func.py"]
