#!/bin/bash
#
# This is a temporary script used to simplify the process of running tests.
# It will be removed at some time point.
#

set -eu
set -o pipefail

rm -rf ./build
mkdir ./build

cp -r *.go ./protos ./build
cd build/protos && ./generate.sh && go mod init polypheny.com/protos && go mod tidy
cd .. && go mod init github.com/polypheny/Polypheny-Go-Driver && go mod edit -replace polypheny.com/protos=./protos  && go mod tidy && clear && go test -v

cd ..

rm -rf ./build
