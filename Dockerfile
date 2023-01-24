FROM alpine:latest

ARG GOLANG_VERSION=1.18.10

#we need the go version installed from apk to bootstrap the custom version built from source
RUN apk update && apk add go gcc bash musl-dev openssl-dev ca-certificates && update-ca-certificates

RUN wget https://dl.google.com/go/go$GOLANG_VERSION.src.tar.gz && tar -C /usr/local -xzf go$GOLANG_VERSION.src.tar.gz

RUN cd /usr/local/go/src && ./make.bash

ENV PATH=$PATH:/usr/local/go/bin

RUN rm go$GOLANG_VERSION.src.tar.gz

#we delete the apk installed version to avoid conflict
RUN apk del go

ENV AES_CRYPT_KEY="n&f+w}IaA^&K\;JhHD>es&Nx7=iH>[gI"
ENV JWT_SECRET="[O?%.%RJHhEGpL&u#Zi-g|b5t:C.m/Kj"

WORKDIR $HOME/gotoko-pos

COPY . .

RUN go build .

CMD ["./gotoko-pos", "-lenv", "false"]