FROM golang:1.19.0-alpine as builder
RUN apk --update add ca-certificates
COPY ./otel-custom /otel-custom
WORKDIR /otel-custom
ENV GO111MODULE=on
ENV CGO_ENABLED=0
RUN go build .


FROM scratch

ARG USER_UID=10001

USER ${USER_UID}

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY  --from=builder /otel-custom/builder /builder
COPY ./config.yml /otelcol-custom/config.yaml
ENTRYPOINT ["/builder"]
CMD ["--config", "/otelcol-custom/config.yaml"]
EXPOSE 4317 55678 55679
