FROM golang:1.17 as build

ARG CGO_ENABLED=0

WORKDIR /app

RUN apt-get update && apt-get install -y upx

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM build as build-sidecar

RUN go build -ldflags="-s -w" -o /cache-server main_container.go \
 && upx --best --lzma /cache-server

FROM build as build-extension

RUN go build -ldflags="-s -w" -o /cache-server main.go \
 && upx --best --lzma /cache-server

FROM scratch

COPY --from=build-sidecar /cache-server /cache-server

EXPOSE 8015

CMD [ "/cache-server" ]