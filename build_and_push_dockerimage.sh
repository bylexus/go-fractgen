#!/bin/bash

HERE=$(dirname "$(readlink -f "$0")")
# if the multiplatform-builder BuildKit instance does not exist yet, create it with:
# docker buildx create --name multiplatform-builder --use
docker buildx use multiplatform-builder
docker buildx build --platform linux/arm64,linux/amd64 -t registry.alexi.ch/go-fractgen:latest --push "${HERE}"

# restart deployment on alexi.ch
ssh kube001.alexi.ch "/snap/bin/microk8s kubectl --namespace alexich rollout restart deployment/fractgen"

