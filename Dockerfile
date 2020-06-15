FROM alpine:latest
RUN adduser -D tidepool
WORKDIR /home/tidepool
USER tidepool
COPY --chown=tidepool ./dist/workscheduler_linux ./workscheduler
CMD ["./workscheduler"]
