FROM golang:1.19.0-alpine as builder
RUN apk --update add ca-certificates
RUN mkdir -p /var/otel
COPY . /var/otel
WORKDIR /var/otel/example/otel-custom
ENV GO111MODULE=on
ENV CGO_ENABLED=0
RUN echo $PWD
RUN go build .


FROM scratch

ARG USER_UID=10001

USER ${USER_UID}

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY  --from=builder /var/otel/example/otel-custom/builder /builder
COPY --from=builder /var/otel/example/config.yml /etc/otelcol-custom/config.yaml
ENTRYPOINT ["/builder"]
CMD ["--config", "/etc/otelcol-custom/config.yaml"]
EXPOSE 4317 55678 55679
