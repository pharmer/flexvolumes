#!/usr/bin/env bash

pushd $GOPATH/src/github.com/pharmer/flexvolumes/hack/gendocs
go run main.go
popd
