FROM alpine:latest

ARG GOLANG_VERSION=1.18.10

#we need the go version installed from apk to bootstrap the custom version built from source
RUN apk update && apk add gcc bash musl-dev openssl-dev ca-certificates && update-ca-certificates

RUN wget https://go.dev/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

# patching up installation
RUN rm go$GOLANG_VERSION.linux-amd64.tar.gz
RUN apk add file patchelf
RUN cd /usr/local/go/bin && patchelf --set-interpreter /lib/libc.musl-x86_64.so.1 go
RUN cd $HOME && go version

ENV AES_CRYPT_KEY="n&f+w}IaA^&K\;JhHD>es&Nx7=iH>[gI"
ENV JWT_SECRET="[O?%.%RJHhEGpL&u#Zi-g|b5t:C.m/Kj"

WORKDIR $HOME/gotoko-pos

COPY . .

RUN go build .

RUN patchelf --set-interpreter /lib/libc.musl-x86_64.so.1 gotoko-pos

RUN apk del -r patchelf file

EXPOSE 3030

CMD ["./gotoko-pos", "-lenv", "false"]