FROM iron/python:2.7.10

WORKDIR /worker
ADD . /worker

ENV PYTHONPATH /worker/packages

ENTRYPOINT ["python", "publish.py"]