FROM golang:alpine
LABEL maintainer "daniel@harrisbaird.co.uk"
RUN apk --update --no-cache add git make ca-certificates
WORKDIR /go/src/github.com/harrisbaird/dailyteedeals
ADD . .
RUN make install && make build && \
    mv /go/src/github.com/harrisbaird/dailyteedeals/bin/dailyteedeals /bin/dailyteedeals
CMD ["/bin/dailyteedeals"]
EXPOSE 8080