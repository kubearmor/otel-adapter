FROM golang:1.22-alpine3.19 as builder

RUN apk --update add ca-certificates git

WORKDIR /build
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0

# install ocb and build
RUN go install go.opentelemetry.io/collector/cmd/builder@latest
RUN builder --config=example/collector-builder.yml

FROM scratch

ARG USER_UID=10001

USER ${USER_UID}

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/otel-custom/otelcol-custom /otelcol

COPY ./example/config.yml /otelcol-custom/config.yml

ENTRYPOINT ["/otelcol"]
CMD ["--config", "/otelcol-custom/config.yml"]
