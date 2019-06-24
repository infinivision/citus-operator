#!/bin/bash

docker build --tag infinivision/stolon-keeper --target keeper . \
&& docker build --tag infinivision/haproxy --target haproxyplus .

# build using proxy
#docker build --network=host --build-arg http_proxy=http://127.0.0.1:8118 --tag infinivision/stolon-keeper --target keeper . \
#&& docker build --network=host --build-arg http_proxy=http://127.0.0.1:8118 --tag infinivision/haproxy --target haproxyplus .
