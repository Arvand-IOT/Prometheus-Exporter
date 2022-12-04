FROM alpine:latest

ARG TARGETOS
ARG TARGETARCH

COPY dist/arvand-exporter_${TARGETOS}_${TARGETARCH} /app/arvand-exporter

COPY scripts/start.sh /app/

EXPOSE 9437

ENTRYPOINT ["/app/start.sh"]
