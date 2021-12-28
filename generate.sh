#!/bin/bash

echo "Generating models..."

rm models/*

docker run --rm --user $(id -u):$(id -g) -v ${PWD}:/local swaggerapi/swagger-codegen-cli-v3 generate \
    -i /local/api/restApi.yaml \
    -l go-server \
    -o /local/models \
    --additional-properties packageName=models \
    --ignore-file-override /local/.swagger-codegen-ignore \
    -D models

mv models/go/* models && rmdir models/go

echo "Generating server..."

docker run --rm --user $(id -u):$(id -g) -v ${PWD}:/local swaggerapi/swagger-codegen-cli-v3 generate \
    -i /local/api/restApi.yaml \
    -l go-server \
    -o /local/tmpserver \
    -t /local/templates/ \
    --additional-properties packageName=server \
    --ignore-file-override /local/.swagger-codegen-ignore

mv tmpserver/go/logger.go tmpserver/go/routers.go server && rm -rf tmpserver
