FROM golang:1.17-alpine as build

ARG CGO_ENABLED=0

RUN apk add upx ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM build as build-sidecar

ARG CGO_ENABLED=0

RUN go build -ldflags="-s -w" -o /cache-server main_container.go \
 && upx --best --lzma /cache-server

FROM build as build-extension

ARG CGO_ENABLED=0

RUN go build -ldflags="-s -w" -o /cache-server main.go \
 && upx --best --lzma /cache-server

FROM scratch

COPY --from=build-sidecar /cache-server /cache-server
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8015

CMD [ "/cache-server" ]