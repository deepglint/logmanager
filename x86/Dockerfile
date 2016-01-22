FROM 192.168.5.46:5000/ubuntu:14.04 

MAINTAINER Kaixian Hu

ADD ./logserver.linux /logserver/logserver.linux

WORKDIR /logserver

CMD ./logserver.linux -influxdb_url 192.168.5.46:8088
