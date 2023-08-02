FROM reg.g5d.dev/node as build
RUN apk add python3 make g++

WORKDIR /src
COPY package.json package-lock.json /src/
RUN npm install

COPY . /src
RUN npm run build
RUN mkdir -p /rootfs/public
RUN cp -r /src/dist/* /rootfs/public/
RUN cp /src/config.toml /rootfs/config.toml

FROM reg.g5d.dev/sws

COPY --from=build /rootfs/ /

ENTRYPOINT ["/sws", "-w/config.toml"]
