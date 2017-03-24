FROM golang:alpine
COPY . $GOPATH/src/app
RUN go build -o /bin/app app
CMD ["/bin/app"]
