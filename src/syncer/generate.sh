#!/bin/bash
rm -rf gen/
#rm -rf cmd/syncer/
rm -f syncer
goa gen syncer/design
goa example syncer/design
export CGO_ENABLED=0
go build -ldflags="-extldflags=-static" ./cmd/syncer
