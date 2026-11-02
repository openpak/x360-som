# Multi-stage, cross-compiled: the build stage runs on the runner's native arch and Go
# emits the target arch, so linux/arm64 builds need no QEMU. No third-party modules.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/x360-som ./cmd/x360-som

FROM alpine:3.22
RUN addgroup -S app && adduser -S -G app app
COPY --from=build /out/x360-som /usr/local/bin/x360-som
COPY modules /srv/modules
ENV X360_SOM_MODULES=/srv/modules
USER app
EXPOSE 10045/tcp
ENTRYPOINT ["/usr/local/bin/x360-som"]
