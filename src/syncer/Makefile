 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=syncer
PREFIX=/usr/local
CMD_DIR=cmd/syncer
DOCKER_PREFIX=../../build
TIMESTAMP=$(shell date +%s)
TAG=build-$(TIMESTAMP)

all: build

build:
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) $(CMD_DIR)/main.go $(CMD_DIR)/http.go $(CMD_DIR)/background.go

goabuild: 
	goa gen syncer/design
	goa example syncer/design

goaclean:
	rm -rf ./gen

install:
	install $(BINARY_NAME) $(PREFIX)/syncer

clean:
	rm -f $(BINARY_NAME)

image:
	rm -f $(DOCKER_PREFIX)/$(BINARY_NAME)
	CGO_ENABLED=0 $(GOBUILD) -o $(DOCKER_PREFIX)/$(BINARY_NAME) $(CMD_DIR)/main.go $(CMD_DIR)/http.go $(CMD_DIR)/background.go
	cd $(DOCKER_PREFIX) \
	&&	docker build . -f Dockerfile -t czdcm-quay.lx.ifortuna.cz/shared-images/syncer:$(TAG) \
	&& 	docker push czdcm-quay.lx.ifortuna.cz/shared-images/syncer:$(TAG)
