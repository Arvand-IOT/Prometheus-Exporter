#!/bin/sh

if [ ! -x /app/arvand-exporter ]; then
  chmod 755 /app/arvand-exporter
fi

if [ -z "$CONFIG_FILE" ]
then
  printf '%s\n' "You should set a config file" >&2
  exit 1
else
    /app/arvand-exporter --config-file $CONFIG_FILE
fi
