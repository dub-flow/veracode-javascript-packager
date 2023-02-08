FROM golang:1.19-alpine

WORKDIR /app

# Download all required dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy all .go files into the container
COPY *.go ./

# Copy the file to compile the app into the container
COPY create-releases.sh ./

# Build the app
RUN ./create-releases.sh docker

# Move into the `/app/js-app` directory where the JavaScript app to packages is
CMD js-app

# Run the Linux x86 release
ENTRYPOINT ["/app/releases/veracode-js-packager-linux-amd64", "-source", ".", "-target", "."]
