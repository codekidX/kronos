FROM golang:alpine AS builder

COPY . .

WORKDIR $GOPATH/src/chrononut/
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/nut

# Create smaller image with only binary
FROM scratch

COPY --from=builder /go/bin/nut /go/bin/nut

ENTRYPOINT ["/go/bin/nut"]
