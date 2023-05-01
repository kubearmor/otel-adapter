FROM alpine:3.16 as certs
RUN apk --update add ca-certificates


FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY  otel-custom/builder /builder
COPY config.yml /etc/otelcol-custom/config.yaml
ENTRYPOINT ["/builder"]
CMD ["--config", "/etc/otelcol-custom/config.yaml"]
EXPOSE 4317 55678 55679
