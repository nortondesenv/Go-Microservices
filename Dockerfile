FROM golang:1.16-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app cmd/main.go

WORKDIR /dist

RUN cp /build/app .

# Build a small image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY . .
COPY --from=builder /dist/app /

EXPOSE 5555

ENTRYPOINT ["/app"] 