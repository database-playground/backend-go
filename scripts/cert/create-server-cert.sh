#!/usr/bin/env bash
# -*- coding: utf-8 -*-

PROPERTY=$1
if [ -z "$PROPERTY" ]; then
  echo "Usage: $0 <property>"
  echo "<property> can be one of 'dev'."
  exit 1
fi

cfssl gencert -ca="ca-$PROPERTY.pem" -ca-key="ca-$PROPERTY-key.pem" -config="ca-config.json" -profile=server "server.${PROPERTY}.json" | cfssljson -bare "server-$PROPERTY"
