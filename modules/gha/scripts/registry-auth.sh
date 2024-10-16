#!/bin/bash

set -o pipefail

echo "$REGISTRY_PASSWORD" | docker login \
  --username "$REGISTRY_USERNAME" --password-stdin  \
  "$REGISTRY_SERVER"
