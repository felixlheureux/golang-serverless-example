FROM golang:1.18

RUN apt-get update
RUN apt-get install python3 python3-pip -y

RUN pip3 install --upgrade pip
RUN pip3 install awscli
RUN pip3 install aws-sam-cli

WORKDIR /var/opt

EXPOSE 8080