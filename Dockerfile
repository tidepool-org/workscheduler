FROM golang:1.14-alpine

WORKDIR /build
COPY . .

RUN set -ex && \
    adduser -D tidepool && \
    mkdir -p /home/tidepool && \
    apk add --no-progress --no-cache gcc musl-dev && \
    go get -d -v ./... && \
    GOOS=linux GOARCH=amd64 go build  -a -v -tags musl -o ./dist/workscheduler ./server && \
    mv ./dist/workscheduler /home/tidepool/workscheduler && \
    chown tidepool /home/tidepool/workscheduler && \
    cd /home/tidepool && \
    rm -rf /build

WORKDIR /home/tidepool
CMD ["./workscheduler"]
