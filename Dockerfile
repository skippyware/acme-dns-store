ARG ALPINE_VERSION=3
ARG GOLANG_VERSION=1.23-alpine

FROM --platform=${BUILDPLATFORM} golang:${GOLANG_VERSION} AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/acme-dns-store


FROM alpine:${ALPINE_VERSION}

RUN apk add --no-cache --no-progress curl
COPY --from=builder /out/acme-dns-store .

CMD ["/acme-dns-store"]

