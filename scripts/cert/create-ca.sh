#!/usr/bin/env bash
# -*- coding: utf-8 -*-

PROPERTY=$1
if [ -z "$PROPERTY" ]; then
  echo "Usage: $0 <property>"
  echo "<property> can be one of 'dev'."
  exit 1
fi

cfssl gencert -initca "ca-csr.$PROPERTY.json" | cfssljson -bare "ca-$PROPERTY" -
