### build go executable
FROM --platform=$BUILDPLATFORM golang:1.21.1 as build

WORKDIR /go/src

COPY go.mod go.sum /go/src/
RUN go mod download

COPY cmd /go/src/cmd
COPY internal /go/src/internal

RUN go test ./...

WORKDIR /go/src/cmd/webhook

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o /go/bin/webhook .

### final image
FROM scratch

ENTRYPOINT ["/app/bin/webhook"]

COPY --from=build /go/bin/webhook /app/bin/webhook
