#!/bin/bash

# Stop and remove the containers
echo -e "Stopping containers...\n"
docker-compose down

# Remove the images
echo -e "Removing images...\n"
docker rmi frontend-image backend-image

echo -e "Cleanup complete."
