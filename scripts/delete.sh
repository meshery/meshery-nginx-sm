#!/bin/sh

set -e

if ! ./scripts/nginx-meshctl remove -y; then
	exit 1
fi
