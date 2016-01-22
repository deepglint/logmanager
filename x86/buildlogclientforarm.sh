#!/bin/bash

TAG=$(<../VERSION)
echo $TAG

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./logserver.linux  ../logserver.go

docker build -t 192.168.5.46:5000/logserver:$TAG .
docker push 192.168.5.46:5000/logserver:$TAG 

rm -rf ./*.linux
