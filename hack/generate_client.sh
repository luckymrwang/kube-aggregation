#!/bin/bash

set -e

GV="$1"

rm -rf ./pkg/client
./hack/generate_group.sh "client,lister,informer" kube-aggregation/pkg/client kubesphere.io/api "${GV}" --output-base=./  -h "$PWD/hack/boilerplate.go.txt"
mv kube-aggregation/pkg/client ./pkg/
rm -rf ./kube-aggregation