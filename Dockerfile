FROM alpine:latest
RUN apk --no-cache add ca-certificates gcompat

WORKDIR /app/

# https://stackoverflow.com/a/35613430
RUN mkdir -p /etc/sentry-prometheus-exporter/

COPY bin/sentry-prometheus-exporter /
COPY config/sentry-prometheus-exporter.yml /etc/sentry-prometheus-exporter/config.yml

EXPOSE 9412
ENTRYPOINT ["/sentry-prometheus-exporter"]
CMD ["--config.file=/etc/sentry-prometheus-exporter/config.yml"]
