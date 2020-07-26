FROM debian:stretch as prepare
RUN DEBIAN_FRONTEND=noninteractive apt-get update 
RUN DEBIAN_FRONTEND=noninteractive apt-get -qq install -y --no-install-recommends ca-certificates gnupg2 curl apt-transport-https
RUN echo "deb https://deb.nodesource.com/node_14.x stretch main" | tee /etc/apt/sources.list.d/nodesource.list
RUN curl -sSL https://deb.nodesource.com/gpgkey/nodesource.gpg.key | apt-key add -
RUN DEBIAN_FRONTEND=noninteractive apt-get update 
RUN DEBIAN_FRONTEND=noninteractive apt-get -qq install -y --no-install-recommends tidy webp bash nodejs jpegoptim optipng
RUN rm -rf /var/lib/apt/lists/*
RUN npm i -g @josee9988/minifyall

WORKDIR /site

COPY favicon.ico /site
COPY index.html /site
COPY assets/ /site/assets/
RUN chown -R nobody:nogroup /site
USER nobody:nogroup

RUN minifyall  -d /site
RUN mv "/site/assets/css/main.css" "/site/assets/css/main-$(sha256sum /site/assets/css/main.css | cut -c -10).css"
RUN sed -i "s#\"assets/css/main.css\"#\"assets/css/$(find /site/assets/css/ -maxdepth 1 -type f -name main-*.css | awk -F/ '{print $NF}')\"#g" /site/index.html
RUN for file in $(find /site -name '*.jpg' -o -name '*.png' -o -name '*.jpeg'); do cwebp -quiet -m 6 -mt -o "$file.webp" -- "$file"; done
RUN jpegoptim -q /site/assets/images/*.jpg
RUN optipng -quiet /site/assets/images/*.png

FROM registry.greboid.com/cv:latest as cv

FROM nginx:mainline-alpine AS nginx

COPY --from=prepare /site /usr/share/nginx/html
COPY --from=cv /srv/http/cv.pdf /usr/share/nginx/html
ADD nginx.conf /etc/nginx/nginx.conf