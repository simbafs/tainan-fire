FROM debian:bullseye-slim as build

WORKDIR /build
LABEL org.opencontainers.image.source="https://github.com/simbafs/coscup-attendance"
LABEL org.opencontainers.image.authors="Simba Fs <me@simbafs.cc>"

ENV PATH="/usr/local/go/bin:/usr/local/node-v18.17.1-linux-x64/bin:$PATH"

# install dependencies
RUN apt-get update && \
    apt-get install -y make xz-utils wget ca-certificates git --no-install-recommends && \
    # config ssl
    mkdir -p /etc/ssl/certs && \
    update-ca-certificates --fresh && \
    # install go v1.23.3
    wget -q https://go.dev/dl/go1.23.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz

COPY . .
RUN go build . -o bot

FROM scratch 

WORKDIR /app

COPY --from=build /build/bot /app/bot
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD [ "./bot" ]
