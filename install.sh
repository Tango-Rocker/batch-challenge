#!/bin/bash

read -p "Enter the directory path for storing demo app data: " data_dir
mkdir -p "$data_dir"

# Set permissions for the data directory (optional, requires sudo)
# read -p "Set permissions for the data directory? [y/N] " permission_choice
# if [[ "$permission_choice" == [Yy] ]]; then
#     sudo chown -R "$USER":"$USER" "$data_dir"
#     sudo chmod -R 755 "$data_dir"
# fi

echo "Building the application..."
./build.sh

# Check if Docker is installed
if ! [ -x "$(command -v docker)" ]; then
    echo "Docker is not installed. Would you like to install it now? (Requires sudo)"
    read -p "Install Docker? [y/N] " install_choice
    if [[ "$install_choice" == [Yy] ]]; then
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
    else
        echo "Docker is required to run this application."
        exit 1
    fi
fi

# Check if Docker Compose is installed
if ! [ -x "$(command -v docker-compose)" ]; then
    echo "Docker Compose is not installed. Would you like to install it now? (Requires sudo)"
    read -p "Install Docker Compose? [y/N] " compose_choice
    if [[ "$compose_choice" == [Yy] ]]; then
        sudo curl -L "https://github.com/docker/compose/releases/download/v2.2.3/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
    else
        echo "Docker Compose is required to run this application."
        exit 1
    fi
fi

# Update the docker-compose.yml file with the chosen data directory
sed -i "s|/path/to/my/data|$data_dir|g" docker-compose.yml

# Run the application using Docker Compose
echo "Starting the application using Docker Compose..."
docker-compose up -d

echo "Installation complete. The application is now running."
