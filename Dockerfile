FROM golang:1.11-alpine
WORKDIR /project
RUN apk add --no-cache git
RUN go get github.com/pdevine/go-asciisprite
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o tetromino *.go

FROM scratch
COPY --from=0 /project/tetromino /tetromino
CMD ["/tetromino"]
