FROM alpine:3.18.3

ARG TARGETOS
ARG TARGETARCH

COPY dist/arvand-exporter_${TARGETOS}_${TARGETARCH} /app/arvand-exporter

COPY scripts/start.sh /app/

RUN chmod +x /app/start.sh

EXPOSE 9437

ENTRYPOINT ["/app/start.sh"]
