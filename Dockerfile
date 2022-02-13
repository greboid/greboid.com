#Generate site with Hugo
FROM reg.g5d.dev/hugo as hugo
COPY --chown=65532:65532 site /build/site
COPY --from=ghcr.io/greboid/cv /cv.pdf /build/site/static/cv.pdf
RUN ["hugo", "-v", "-s", "/build/site"]

#Minify + Image optimisation
FROM reg.g5d.dev/alpine as minify
RUN apk add --no-cache libwebp-tools;
COPY --from=hugo --chown=65532:65532 /build/public /tmp/public
USER 65532:65532
RUN find /tmp/public \( -name '*.jpg' -o -name '*.png' -o -name '*.jpeg' \) -exec cwebp -q 60 "{}" -o "{}.webp" \;;

#Serve with nginx
FROM reg.g5d.dev/nginx:latest AS nginx
COPY --from=minify /tmp/public /src/public
COPY nginx.conf /usr/local/nginx/conf/nginx.conf

ENTRYPOINT ["/nginx"]
