#!/usr/bin/env sh

# get the current version of the tool from `./current_version`
VERSION=$(cat current_version)

docker buildx build --platform linux/amd64,linux/arm64 --push . -t fw10/veracode-js-packager:latest -t fw10/veracode-js-packager:$VERSION