FROM debian:stretch-20220509-slim

EXPOSE 9437

COPY scripts/start.sh /app/

COPY dist/arvand-exporter_linux_amd64 /app/arvand-exporter

RUN chmod 755 /app/*

ENTRYPOINT ["/app/start.sh"]
