#!/usr/bin/env bash

pushd $GOPATH/src/github.com/pharmer/pharmer/hack/gendocs
go run main.go
popd
