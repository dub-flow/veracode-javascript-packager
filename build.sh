# build for Windows
GOOS=windows GOARCH=amd64 go build
mv vc-node-packager.exe releases/windows-amd64

# build for Mac M1 (arm64)
GOOS=darwin GOARCH=arm64 go build
mv vc-node-packager releases/mac-arm64

# build for 64bit Linux (amd64)
GOOS=linux GOARCH=amd64 go build
mv vc-node-packager releases/linux-amd64
