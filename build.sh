#! /bin/sh

docker build . -t formable || exit 1
docker run \
    -v `pwd`/example/:/mnt/ \
    -e TF_CMD=plan \
    -e PROJECT=dev-cicd \
    -e SERVICE_ACCOUNT=omni-king \
    -e SERVICE=example \
    -e REGION=us-central1 \
    -e ZONE=us-central1-a \
    formable:latest