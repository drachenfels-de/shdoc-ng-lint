FROM golang:1.25.1-alpine AS build

WORKDIR /src
ARG TARGETOS
ARG TARGETARCH
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -trimpath -ldflags="-s -w" -o /out/shdoc-ng-lint ./cmd/shdoc-ng-lint

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=build /out/shdoc-ng-lint /usr/local/bin/shdoc-ng-lint
ENTRYPOINT ["shdoc-ng-lint"]
