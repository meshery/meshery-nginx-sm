#!/bin/sh

set -e

if ! scripts/nginx-meshctl deploy --nginx-mesh-api-image "${NGINX_DOCKER_REGISTRY}/nginx-mesh-api:${NGINX_MESH_VER}" --nginx-mesh-sidecar-image "${DOCKER_REGISTRY}/nginx-mesh-sidecar:${NGINX_MESH_VER}" --nginx-mesh-init-image "${DOCKER_REGISTRY}/nginx-mesh-init:${NGINX_MESH_VER}" --nginx-mesh-metrics-image "${DOCKER_REGISTRY}/nginx-mesh-metrics:${NGINX_MESH_VER}"; then
	exit 1
fi
