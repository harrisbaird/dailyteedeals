# Build binary
FROM golang:alpine AS build-env
LABEL maintainer "daniel@harrisbaird.co.uk"

RUN apk --update --no-cache add git ca-certificates
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/harrisbaird/dailyteedeals
ADD . .
RUN dep ensure
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /bin/dailyteedeals

# Small runtime image
FROM alpine
RUN apk --update --no-cache add ca-certificates
COPY --from=build-env /bin/dailyteedeals /bin/dailyteedeals
ENTRYPOINT /bin/dailyteedeals