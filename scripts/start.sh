#!/bin/sh

if [ ! -x /app/arvand-exporter ]; then
  chmod 755 /app/arvand-exporter
fi

if [ -z "$CONFIG_FILE" ]
then
    /app/arvand-exporter -device $DEVICE -address $ADDRESS -user $USER -password $PASSWORD
else
    /app/arvand-exporter -config-file $CONFIG_FILE
fi
