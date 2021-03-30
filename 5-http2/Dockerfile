FROM golang:1.13
WORKDIR /src
ADD server.go server.go
RUN go build server.go
CMD ["./server"]
EXPOSE 8080
