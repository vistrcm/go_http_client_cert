#!/usr/bin/env bash
set -xv

curl -v \
  --cacert tls/server_cert.pem \
  --cert tls/client_cert.pem \
  --key tls/client_key.pem \
  https://localhost:8443/hello
