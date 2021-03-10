FROM registry.greboid.com/mirror/debian:latest as webp
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get -qq install -y --no-install-recommends webp python && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY static/ /app/static/
COPY minify.sh /app
RUN /bin/bash /app/minify.sh

FROM registry.greboid.com/cv:latest as cv

FROM registry.greboid.com/mirror/golang:latest as builder
WORKDIR /app
COPY --from=cv /srv/http/cv.pdf /app/static/cv.pdf
COPY --from=webp /app/static/ /app/static/
COPY main.go /app
COPY handlers.go /app
COPY go.mod /app
COPY go.sum /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main .

FROM scratch
WORKDIR /app
COPY --from=builder /app/main /app
EXPOSE 8080
CMD ["/app/main"]