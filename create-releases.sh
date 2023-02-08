#!/usr/bin/env sh

VERSION="1.0.0"

FLAGS="-X main.AppVersion=$VERSION -s -w"

rm -rf releases
mkdir -p releases

# check if `./create-releases.sh docker` is ran which would only trigger the x64 Linux release being built
if [[ $1 != "docker" ]]; then
    # build for Windows
    GOOS=windows GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath
    mv veracode-js-packager.exe releases/veracode-js-packager-windows-amd64.exe

    # build for M1 Macs (arm64)
    GOOS=darwin GOARCH=arm64 go build -ldflags="$FLAGS" -trimpath
    mv veracode-js-packager releases/veracode-js-packager-mac-arm64

    # build for Intel Macs (amd64)
    GOOS=darwin GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath
    mv veracode-js-packager releases/veracode-js-packager-mac-amd64
fi

# build for x64 Linux (amd64)
GOOS=linux GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath
mv veracode-js-packager releases/veracode-js-packager-linux-amd64
