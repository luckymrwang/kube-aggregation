FROM golang:1.19.5

# FROM golang:1.19.2-alpine
# RUN apk add --no-cache build-base git bash

RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu gcc-x86-64-linux-gnu

COPY . /kube-aggregation

# RUN rm -rf /kube-aggregation/.git
# RUN rm -rf /kube-aggregation/test

ENV CLUSTERPEDIA_REPO="/kube-aggregation"

RUN cp /kube-aggregation/hack/builder.sh /
