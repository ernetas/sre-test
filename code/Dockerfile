FROM golang:latest as builder
ADD ./ /app
WORKDIR /app
RUN go get -d -v && \
    go build
EXPOSE 9000
CMD ["/app"]
