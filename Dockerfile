FROM golang:1.25.1-alpine AS build

WORKDIR /work
ARG TARGETOS
ARG TARGETARCH
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} make build

FROM scratch
COPY --from=build /work/shdoc-ng-lint /shdoc-ng-lint
ENTRYPOINT ["/shdoc-ng-lint"]
