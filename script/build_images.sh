#!/usr/bin/env bash
#you need to login docker first
#docker login 192.168.0.251:5005
#use another domain "peakperformance.ddns.net:251" to pull image
image=192.168.0.251:5005/supersign/crypto-wallet/server:v1.0.1
cd ../
docker build -t  $image . -f ./docker/images/share-wallet/deploy.Dockerfile
docker push $image
