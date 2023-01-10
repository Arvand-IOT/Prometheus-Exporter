# Arvand Prometheus Exporter

![Go](https://github.com/Arvand-IOT/Prometheus-Exporter/workflows/Go/badge.svg)

![banner](banner.jpg)

```bash
docker run -itd -n arvand-exporter -p 9437:9437 -v <path-to-config.yml>:/app/config.yml -e CONFIG_FILE="/app/config.yml" hatamiarash7/arvand-exporter:1.2.0
```

Sample `config.yml` :

```yml
clients:
  - name: Room 1
    ip: 192.168.1.2
  - name: Room 2
    ip: 192.168.1.3
  - name: Room 3
    ip: 192.168.1.4
  - name: Room 4
    ip: 192.168.1.5
```

## Metrics

We have these metrics at this moment:

- Temperature
- Humidity
- Air Quality
- Light
