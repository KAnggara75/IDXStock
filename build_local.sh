#!/bin/bash

# Build the image
echo "Building podman image: idxstock:latest..."
podman build -f Dockerfile -t idxstock:latest .

# Apply the kube manifests
echo "Applying podman kube manifests..."
podman play kube --replace build/local.yaml

# Cleanup dangling images
echo "Cleaning up dangling images..."
DANGLING_IMAGES=$(podman images --filter "dangling=true" --quiet)
if [ -n "$DANGLING_IMAGES" ]; then
    podman rmi -f $DANGLING_IMAGES
    echo "Removed dangling images."
else
    echo "No dangling images found."
fi

echo "Process completed successfully."
