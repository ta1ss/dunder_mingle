# !/bin/bash

# Build the images
echo -e "Building images...\n"
docker build -t backend-image .

cd front-end
docker build -t frontend-image .
cd ..

# Wait for the images to be built
while ! docker images | grep frontend-image && ! docker images | grep backend-image; do
  echo -e "Something might be wrong? \n"
  sleep 1
done

# Run the containers
echo -e "Starting images...\n"

docker-compose up