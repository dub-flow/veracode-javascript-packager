# This is the docker image we use to build the app
FROM golang:1.19-alpine as build

WORKDIR /app

# Download all required dependencies
COPY go.mod ./
COPY go.sum ./

# Copy all .go files into the container
COPY *.go ./
# Copy the file to compile the app into the container
COPY create-releases.sh ./
COPY current_version ./

# get e.g. `gcc` to compile the app (because this is not part of the `golang:*-alpine` images)
RUN apk add build-base

# Build the app
RUN ./create-releases.sh docker

# This is the much smaller docker image which we will use to run the app
FROM alpine

# Running as a non-root user
RUN adduser -D local
USER local

# Copy the compiled app into the distroless image
COPY --from=build /app /app

# Move into the `/app/js-app` directory where the JavaScript app to package is
WORKDIR /app/js-app

# Run the Linux x86 release
ENTRYPOINT ["/app/veracode-js-packager", "-source", ".", "-target", "."]