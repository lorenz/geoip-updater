FROM scratch

ADD cacert.pem /etc/ssl/ca-bundle.pem

SUB golang:1.7
ENV CGO_ENABLED 0
ADD . /go/src/geoip-updater
WORKDIR /go/src/geoip-updater
RUN go get .
RUN go build .
RETURN /go/src/geoip-updater/geoip-updater /geoip-updater

ENV USER_ID 99999
ENV EDITION_IDS GeoLite2-City

VOLUME /data
CMD ["/geoip-updater"]
