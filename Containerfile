FROM golang:1.25.5-alpine as builder

WORKDIR /app
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 go build -trimpath -tags 'netgo,osusergo' -ldflags='-s -w -extldflags "-static"' -o main .

FROM ghcr.io/greboid/dockerbase/nonroot:1.20251213.0
COPY --from=builder /app/main /greboid.com
EXPOSE 8080
ENTRYPOINT ["/greboid.com"]
