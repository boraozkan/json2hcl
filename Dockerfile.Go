# syntax = docker/dockerfile:1.2

FROM --platform=${BUILDPLATFORM} golang:1.16-alpine AS base
WORKDIR /src
ENV CGO_ENABLED=0

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS build
ARG TARGETOS
ARG TARGETARCH

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} cd json2hcl && go build -o /out/json2hcl .

FROM alpine AS bin-unix
COPY --from=build /out/json2hcl .

FROM bin-unix AS bin-linux
FROM bin-unix AS bin-darwin

FROM bin-${TARGETOS} as bin

CMD [ "./json2hcl" ]