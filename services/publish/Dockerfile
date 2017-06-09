FROM python:2.7.13-alpine

WORKDIR /worker
ADD . /worker

ENV PYTHONPATH /worker/packages

ENTRYPOINT ["python", "publish.py"]