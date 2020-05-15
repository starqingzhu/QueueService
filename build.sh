#!/bin/bash

mkdir -p bin

rm -rf bin

go build -o bin/client client/client.go
go build -o bin/server server/server.go