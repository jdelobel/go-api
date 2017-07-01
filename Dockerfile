FROM alpine

ADD ./config/config.json /config/config.json

ADD ./apid /

ENTRYPOINT ["/apid"]
