FROM --platform=linux/amd64 golang:1.21 as build

WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o app cmd/moscode-cdn-fiber/main.go

FROM scratch
COPY --from=build ["/build/app", "/"]

EXPOSE 8080

CMD ["/app"]
