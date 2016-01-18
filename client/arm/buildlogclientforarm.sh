#!/bin/bash

TAG=$(<../VERSION)
echo $TAG

GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o ./logclient.arm  ../logclient.go

docker build -t 192.168.5.46:5000/armhf-logclient:$TAG .
docker push 192.168.5.46:5000/armhf-logclient:$TAG 

rm -rf ./*.linux
