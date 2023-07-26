FROM golang:1.20-alpine3.17 as builder

RUN apk --update add ca-certificates git

WORKDIR /build
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0

# install ocb and build
RUN go install go.opentelemetry.io/collector/cmd/builder@latest
RUN builder --config=example/collector-builder.yml

FROM alpine

ARG USER_UID=10001

USER ${USER_UID}

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/otel-custom/otelcol-custom /builder

COPY ./example/config.yml /otelcol-custom/config.yaml

EXPOSE 4317 55678 55679

ENTRYPOINT ["/builder"]
CMD ["--config", "/otelcol-custom/config.yaml"]
