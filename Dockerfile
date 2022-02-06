FROM reg.g5d.dev/alpine:latest AS images

COPY . /src

RUN set -eux; \
    apk add libwebp-tools git coreutils bash; \
    cd /src; \
    mkdir -p /app/static/images; \
    cp -r static/images /app/static; \
    chmod +x minify.sh; \
    ./minify.sh

FROM ghcr.io/greboid/cv:latest as cv

FROM reg.g5d.dev/golang:latest AS build

COPY --from=images /app/static/images /images
COPY --from=cv /cv.pdf /cv.pdf
COPY --from=cv /reversed.pdf /reversed.pdf

COPY . /src

RUN set -eux; \
    apk add git; \
    cd /src; \
    cp -r /images/* static/images/; \
    cp /cv.pdf static/cv.pdf; \
    cp /reversed.pdf static/reversed.pdf; \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -trimpath -ldflags=-buildid= -o main .; \
    # Clobber all the timestamps to make the build more reproducible
    touch --date=@0 /src;

FROM reg.g5d.dev/base:latest

COPY --from=build /src/main /greboid.com
EXPOSE 8080
ENTRYPOINT ["/greboid.com"]