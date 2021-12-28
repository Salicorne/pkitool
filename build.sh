#!/bin/bash

# sudo mkdir /sys/fs/cgroup/systemd
# sudo mount -t cgroup -o none,name=systemd cgroup /sys/fs/cgroup/systemd

docker run --rm -v $(pwd):/workspace -w /workspace --network=host golang:1.16 go build