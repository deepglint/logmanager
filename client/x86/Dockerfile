FROM 192.168.5.46:5000/ubuntu:14.04 

MAINTAINER Kaixian Hu

ADD ./logclient.linux /logclient/logclient.linux

WORKDIR /logclient

CMD ./logclient.linux -keep_interval 30m -upload_interval 10m -sleep_interval 10m -server_host http://192.168.5.46 -dir /tmp/
