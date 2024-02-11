# HOODID API

### Go Base Image
FROM golang:1.20.4-alpine3.17 AS base_builder

WORKDIR /build/

COPY ["go.mod", "go.sum", "./"]

RUN go mod download


### Build Go
FROM base_builder AS go_build

WORKDIR /build/

COPY . .

ARG PROJECT_VERSION=1 CI_COMMIT_SHORT_SHA=1
RUN go build -ldflags="-s -w -X 'main.VERSION=$PROJECT_VERSION' -X main.COMMIT=$CI_COMMIT_SHORT_SHA" -o app cmd/api/main.go


### Build Docker Image
FROM alpine:3.17

WORKDIR /app/

COPY --from=go_build ["/build/app", "./"]

EXPOSE 8080

ENTRYPOINT ["./app"]

#export PROJECT_VERSION=$(cat $(pwd)/VERSION)
#export CI_COMMIT_SHORT_SHA=$(git rev-parse --short HEAD) ou pegar a $CI_COMMIT_SHORT_SHA do gitlab
#docker build --build-arg PROJECT_VERSION=$(cat $(pwd)/VERSION) --build-arg CI_COMMIT_SHORT_SHA=$(git rev-parse --short HEAD) -t acragentesvirtuaisdev.azurecr.io/plataforma-go-security-gateway:$(cat $(pwd)/VERSION) -t acragentesvirtuaisdev.azurecr.io/plataforma-go-security-gateway:latest .