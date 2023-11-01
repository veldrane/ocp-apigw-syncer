#!/bin/bash
rm -rf gen/
#rm -rf cmd/syncer/
rm -f syncer
goa gen syncer/design
goa example syncer/design
go build ./cmd/syncer
