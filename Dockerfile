FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
WORKDIR /go/gcs-uploader
COPY go.* ./
RUN go mod download
COPY . .
ARG TARGETARCH
ARG TARGETOS

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -ldflags "-s -w" -o /out/gcs-uploader .

FROM alpine:3
COPY --from=build /out/gcs-uploader /usr/bin/
ENTRYPOINT ["gcs-uploader"]
