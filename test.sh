#!/bin/bash
set -e

host=127.0.0.1
port=8080
addr="http://$host:$port"

pkiname="my-pki"

# Create PKI
echo ""
echo "Creating PKI"
echo ""
curl -s -XPOST $addr/pki \
  -H 'Content-Type: application/json' \
  -d "{\"name\":\"$pkiname\",\"validityDays\":3650,\"DN\":{\"country\":[\"FR\"],\"locality\":[\"Rennes\"],\"organization\":[\"Salicorne\"],\"commonName\":\"$pkiname\"}}"

# Get full PKI
echo ""
echo "Getting full PKI"
echo ""
curl -s -XGET $addr/pki/$pkiname | jq .
subCaRequestUrl=$(curl -s -XGET $addr/pki/$pkiname | jq -r .subCaRequestUrl)

# Create first subCA
echo ""
echo "Creating subCA"
echo ""
curl -s -XPOST $addr$subCaRequestUrl \
  -H 'Content-Type: application/json' \
  -d "{\"validityDays\":365,\"DN\":{\"country\":[\"FR\"],\"locality\":[\"Rennes\"],\"organization\":[\"Salicorne\"],\"commonName\":\"$pkiname-subca-1\"}}"

# Get full PKI
echo ""
echo "Getting full PKI"
echo ""
curl -s -XGET $addr/pki/$pkiname | jq .

# Create second subCA
echo ""
echo "Creating subCA"
echo ""
curl -s -XPOST $addr$subCaRequestUrl \
  -H 'Content-Type: application/json' \
  -d "{\"validityDays\":365,\"DN\":{\"country\":[\"FR\"],\"locality\":[\"Rennes\"],\"organization\":[\"Salicorne\"],\"commonName\":\"$pkiname-subca-2\"}}"

# Get full PKI
echo ""
echo "Getting full PKI"
echo ""
curl -s -XGET $addr/pki/$pkiname | jq .