# Syncer

## Vision

The aim of this project is checking consistency of distributed key/val store in nginx plus in kubernetes environment.


## Installation

We’re using these technologies: golang, goa framework

Once installed, you have to run these commands to setup the project (example for cz test):

```shell
make install
cd ../build
oc project apigwp-cz
helm install ng-apigw-syncer . --values=files/values/apigwp/cz/t/values.yaml --set ocp4=true
```

Nginx apigw must be configured to use this tool


## How it works

Syncer has a two main part - scraper process and main http handler

scraper process - periodically check openshift api and set internal list of the pods/nginx instancies.

http handler - checks the presence of the token based on the replicas


## How http handler works 

Syncer get the resp api request on /v1/synced with two inputs:
    - token (header X-Auth-Token)
    - origin pod (header X-Nginx-Origin)

Handler start multiple gorutine in single waitgrou for each instance of the nginx except the origin pod. The list of the pods is prepared by scraper process.
Each gorutine tries to connect to own nginx instance on the jwt secured endpoint. If the response from the nginx instance is 200, then gorutine knows that this
instance has token so its synced and finishes. If the response is 401 then, token is not present here. In this case, gorutine waits for some short delay and tries
the request again until gets 200, reach the limits of the retries or reach the timeout of the context (hard timeout). After all gorutines finish, code evaluate
the final state of each gorutine and returns one of the four state in header X-Token-Status:

    - Synced (all gorutines finishes with 200 code in the end)
    - Partialy (more than 50% gorutines finishes with 200 code in the end)
    - NotSynced (less than 50% gorutines finishes with 200 code in the end)
    - Timeout (the all request were not finished until context deadline)


## Project directory content

./build - Dockerfile a all stuf necessary for syncer building
./helm - helm chart
./roles - syncer needs to have access for listing pods and replication sets inside the namespace. Role inside this directory contains these privileges. On new cluster needs to be installed first
./slides - reveal-md presentation, css, images and all stuff
./src - golang sources code

## Important source directory and files in ./src/syncer

./design - design files for goa framework - look at make build

./local/ocp4cli - library for communication with openshift 4.x cluster

./local/synclib - library contains core logic of http handler

./gen - do not touch! - generated directory by goa framework. Customize code in this directory is not recommended

./public - swagger client and openapi.json definition

./cmd/syncer/main.go - entry point for syncer tool

./cmd/syncer/background.go - code of the scraper process