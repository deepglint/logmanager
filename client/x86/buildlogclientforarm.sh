#!/bin/bash

TAG=$(<../VERSION)
echo $TAG

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./logclient.linux  ../logclient.go

docker build -t 192.168.5.46:5000/logclient:$TAG .
docker push 192.168.5.46:5000/logclient:$TAG 

rm -rf ./*.linux
