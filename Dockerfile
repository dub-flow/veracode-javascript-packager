# This is the docker image we use to build the app
FROM golang:1.19-alpine as build

RUN apk add build-base

WORKDIR /app

# Download all required dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy all .go files into the container
COPY *.go ./

# Copy the file to compile the app into the container
COPY create-releases.sh ./

# Build the app
RUN ./create-releases.sh docker

# This is the much smaller docker image which we will use to run the app
FROM alpine:latest

WORKDIR /app
# Copy the compiled app into the distroless image
COPY --from=build /app .

# Move into the `/app/js-app` directory where the JavaScript app to packages is
WORKDIR /app/js-app

# Run the Linux x86 release
ENTRYPOINT ["/app/veracode-js-packager", "-source", ".", "-target", "."]