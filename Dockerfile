#Reverse page
FROM reg.g5d.dev/alpine as reversed
ARG REVERSE=false
RUN apk add patch
COPY reversed.patch /build/reversed.patch
COPY site/themes/greboid.com/assets/sass/main.scss /build/site/themes/greboid.com/assets/sass/main.scss
COPY site/config.yaml /build/site/config.yaml
COPY --from=ghcr.io/greboid/livecv /cv.pdf /cv.pdf
COPY --from=ghcr.io/greboid/livecv /reversed.pdf /reversed.pdf
RUN  if [ "${REVERSE}" == "true" ]; then patch /build/site/themes/greboid.com/assets/sass/main.scss < /build/reversed.patch; fi
RUN  if [ "${REVERSE}" == "true" ]; then mv /reversed.pdf /cv.pdf; fi
RUN  if [ "${REVERSE}" == "true" ]; then sed -i 's+https://greboid.com/+https://dioberg.co.uk+' /build/site/config.yaml; fi

#Generate site with Hugo
FROM reg.g5d.dev/hugo as hugo
COPY --chown=65532:65532 site /build/site
COPY --chown=65532:65532 --from=reversed /build/site/config.yaml /build/site/config.yaml
COPY --chown=65532:65532 --from=reversed /build/site/themes/greboid.com/assets/sass/main.scss /build/site/themes/greboid.com/assets/sass/main.scss
COPY --from=reversed /cv.pdf /build/site/static/cv.pdf
RUN ["hugo", "-v", "-s", "/build/site"]

#Minify + Image optimisation
FROM reg.g5d.dev/alpine as minify
RUN apk add --no-cache libwebp-tools brotli libavif-apps;
COPY --from=hugo --chown=65532:65532 /build/public /tmp/public
USER 65532:65532
RUN find /tmp/public \( -name '*.jpg' -o -name '*.png' -o -name '*.jpeg' \) -exec cwebp -q 60 "{}" -o "{}.webp" \; -exec avifenc -j all --max 40 --maxalpha 63 -r l -s 0 "{}" "{}.avif" \; ;\
    find /tmp/public -name *.jpg -exec mv "{}" "{}.jpg" \; ;\
    find /tmp/public -name *.jpeg -exec mv "{}" "{}.jpeg" \; ;\
    find /tmp/public -name *.png -exec mv "{}" "{}.png" \; ;\
    find /tmp/public \( -name '*.css' -o -name '*.html' -o -name '*.xml' \) -exec brotli --best "{}" \; -exec gzip --best -k {} \; ; \
    find /tmp/public -name '*.css' -exec mv "{}" "{}.css" \; ;\
    find /tmp/public -name '*.html' -exec mv "{}" "{}.html" \; ;\
    find /tmp/public -name '*.xml' -exec mv "{}" "{}.xml" \;

#Serve with nginx
FROM reg.g5d.dev/apache
COPY --chown=65532:65532 --from=minify /tmp/public /usr/local/apache2/htdocs
COPY httpd.conf /usr/local/apache2/conf/httpd.conf
