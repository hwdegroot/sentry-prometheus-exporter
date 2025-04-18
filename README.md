# Sentry exporter

An exporter for [Prometheus](https://prometheus.io/) that collects metrics from [Sentry](https://sentry.io).

## Install

You can download pre-built binaries from [GitHub releases](https://github.com/snakecharmer/sentry_exporter/releases)

## Building and running

### Local Build

```sh
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
go get -v
go build -o bin/sentry-prometheus-exporter src/*
./bin/sentry_exporter <flags>
```

or you can use gox:

```sh
go get github.com/mitchellh/gox
gox
```

Visiting [http://localhost:9412/probe?target={sentry_project}](http://localhost:9412/probe?target=google.com)
will return metrics for a probe against the sentry project

### Building with Docker

  ```sh
  docker build -t sentry-prometheus-exporter .
  docker run -d \
    -p 9412:9412 \
    --name sentry-prometheus-exporter \
    -v `pwd`:/config sentry-prometheus-exporter \
    --config.file=/config/sentry-prometheus-exporter.yml
  ```

## [Configuration](CONFIGURATION.md)

Sentry exporter is configured via a [configuration file](CONFIGURATION.md) and
command-line flags (such as what configuration file to load, what port to listen
on, and the logging format and level).

Sentry exporter can reload its configuration file at runtime. If the new configuration
is not well-formed, the changes will not be applied.
A configuration reload is triggered by sending a `SIGHUP` to the Sentry exporter
process or by sending a HTTP POST request to the `/-/reload` endpoint.

To view all available command-line flags, run `./sentry_exporter -h`.

To specify which [configuration file](CONFIGURATION.md) to load, use the `--config.file`
flag.

Additionally, an [example configuration](sentry_exporter.yml) is also available.

The timeout of each probe is automatically determined from the `scrape_timeout`
in the [Prometheus configuration file](https://prometheus.io/docs/operating/configuration/#configuration-file),
slightly reduced to allow for network delays.

## Prometheus Configuration

The sentry exporter needs to be passed the target as a parameter, this can be
done with relabelling.

Example configuration:

```yml
scrape_configs:
  - job_name: 'sentry'
    metrics_path: /probe
    static_configs:
      - targets:
        - {project name}    # First project name.
        - {second project name}   # Second project name
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9412  # The sentry exporter's real hostname:port.
```
