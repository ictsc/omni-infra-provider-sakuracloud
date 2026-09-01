FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' \
    -o /out/omni-infra-provider-sakuracloud \
    ./cmd/omni-infra-provider-sakuracloud

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/omni-infra-provider-sakuracloud /omni-infra-provider-sakuracloud
ENTRYPOINT ["/omni-infra-provider-sakuracloud"]
