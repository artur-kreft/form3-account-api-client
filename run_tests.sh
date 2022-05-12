#!/bin/sh
go test ./... -cover -v -args -api_url $1
exit 1
