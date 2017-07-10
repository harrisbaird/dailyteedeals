# Build binary
FROM golang:alpine AS build-env
LABEL maintainer "daniel@harrisbaird.co.uk"

RUN apk --update --no-cache add git ca-certificates

# install dependency tool
RUN go get github.com/golang/dep && go install github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/harrisbaird/dailyteedeals
ADD . .
RUN dep ensure
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/dailyteedeals && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o migrate migrations/*.go

# Small runtime image
FROM alpine
RUN apk --update --no-cache add ca-certificates
COPY --from=build-env /bin/dailyteedeals /bin/dailyteedeals
COPY --from=build-env /bin/migrate /bin/migrate
COPY entrypoint.sh /bin/entrypoint.sh
ENTRYPOINT /bin/entrypoint.sh
EXPOSE 8080