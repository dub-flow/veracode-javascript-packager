#!/usr/bin/env sh

docker buildx build --platform linux/amd64,linux/arm64 --push . -t fw10/veracode-js-packager:latest -t fw10/veracode-js-packager:1.0.0