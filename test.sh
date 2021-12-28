#!/bin/bash
set -e

docker run --rm -v $(pwd):/workspace -w /workspace --network=host golang:1.16 go build

./pkitool