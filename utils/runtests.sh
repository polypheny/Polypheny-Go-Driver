#!/bin/bash
#
# This is a temporary script used to simplify the process of running tests.
# It will be removed at some time point.
#

set -eu
set -o pipefail

cd ..

rm -rf ./build
mkdir ./build

cp -r ./protos ./utils/generate.go ./build/

cd build && go run generate.go && rm -rf protos generate.go && cd ..

cp *.go ./build

cd build

cd polypheny.com/protos/ && go mod init polypheny.com/protos && go mod tidy && cd ../..

go mod init github.com/polypheny/Polypheny-Go-Driver && go mod edit -replace polypheny.com/protos=./polypheny.com/protos  && go mod tidy && clear && go test -v

cd ..

rm -rf ./build
