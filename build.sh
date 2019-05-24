#!/bin/bash

sudo docker build -t infinivision/stolon-keeper -f Dockerfile_keeper .

sudo docker build -t infinivision/stolon-proxy -f Dockerfile_proxy .

sudo docker build -t infinivision/stolon-sentinel -f Dockerfile_sentinel .
