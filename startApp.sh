#!/bin/bash
podman stop injectionapp
podman container rm injectionapp
podman image rm injectionapp

podman build --tag injectionapp:latest .
podman run --name injectionapp --network="host" injectionapp
