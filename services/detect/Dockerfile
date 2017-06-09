FROM denismakogon/for-chad:latest

WORKDIR /func
ADD . /func

CMD sh -c 'ln /dev/null /dev/raw1394'; TASK_ID=123 PAYLOAD_FILE=/worker/payload_direct.json python func.py

EXPOSE 8080
