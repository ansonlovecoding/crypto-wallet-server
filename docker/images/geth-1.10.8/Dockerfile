FROM ethereum/client-go:v1.10.26

RUN apk update && apk add bash curl

ADD ./script /root/script
RUN chmod a+x /root/script/*

ENTRYPOINT /root/script/start.sh