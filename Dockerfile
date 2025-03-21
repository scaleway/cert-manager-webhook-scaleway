FROM golang:1.24.1-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /go/src/github.com/arbreagile/cert-manager-webhook-bunny

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=readonly -a -o cert-manager-webhook-bunny main.go

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/arbreagile/cert-manager-webhook-bunny/cert-manager-webhook-bunny .
ENTRYPOINT ["/cert-manager-webhook-bunny"]