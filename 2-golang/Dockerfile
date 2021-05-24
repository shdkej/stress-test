FROM golang:1.14-alpine
WORKDIR /src
ADD server.go server.go
RUN go build server.go
EXPOSE 8080
CMD ["./server"]
