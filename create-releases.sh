rm -rf releases
mkdir -p releases

# build for Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath
mv veracode-js-packager.exe releases/veracode-js-packager-windows-amd64.exe

# build for Mac M1 (arm64)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath
mv veracode-js-packager releases/veracode-js-packager-mac-arm64

# build for x64 Linux (amd64)
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath
mv veracode-js-packager releases/veracode-js-packager-linux-amd64
