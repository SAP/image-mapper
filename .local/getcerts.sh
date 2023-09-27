#!/usr/bin/env bash

set -eo pipefail

cd "$(dirname "$0")"

mkdir -p ssl

kubectl get secret image-mapper-webhook -o jsonpath='{.data.tls\.key}' | base64 -d > ssl/tls.key
kubectl get secret image-mapper-webhook -o jsonpath='{.data.tls\.crt}' | base64 -d > ssl/tls.crt