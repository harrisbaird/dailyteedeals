# Build binary
FROM golang:alpine AS build-env
LABEL maintainer "daniel@harrisbaird.co.uk"
RUN apk --update --no-cache add git ca-certificates make
WORKDIR /go/src/github.com/harrisbaird/dailyteedeals
ADD . .
RUN make install_dependencies && \
    BUILD_LOCATION="/bin" make build

# Small runtime image
FROM alpine
RUN apk --update --no-cache add ca-certificates
COPY --from=build-env /bin/dailyteedeals /bin/dailyteedeals
COPY --from=build-env /bin/migrate /bin/migrate
COPY entrypoint.sh /bin/entrypoint.sh
ENTRYPOINT /bin/entrypoint.sh
EXPOSE 8080