#!/bin/bash

set -o pipefail

if [[ ! -f "$HOME/.docker/config.json" ]]
then
  echo "$DOCKERHUB_PASSWORD" | docker login \
    --username "$DOCKERHUB_USERNAME" --password-stdin  \
    "$DOCKERHUB_SERVER"
fi
