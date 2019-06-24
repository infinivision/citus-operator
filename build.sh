#!/bin/bash

docker build --tag infinivision/stolon-keeper --target keeper . \
&& docker build --tag infinivision/haproxy --target haproxyplus .
