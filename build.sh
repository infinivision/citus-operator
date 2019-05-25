#!/bin/bash

docker build --tag infinivision/stolon-keeper --target keeper .
docker build --tag infinivision/stolon-proxy --target proxy .
docker build --tag infinivision/stolon-sentinel --target sentinel .
