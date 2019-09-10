FROM ubuntu:18.04

RUN apt-get update -y && \
      apt-get install -y python-pip python-dev python3-distutils 


COPY ./requirements.txt /app/requirements.txt

WORKDIR /app

RUN pip install -r requirements.txt

COPY . /app


ENV FLASK_APP node_server.py


CMD [ "flask", "run", "--host", "0.0.0.0" ]
