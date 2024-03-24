#!/bin/bash
echo "Stopping and then deleting existing DB container"
podman stop db
podman container rm -f db

echo "Clearing out any existing DB data files"
sudo rm -rf /tmp/data
mkdir /tmp/data

echo "Starting mysql DB on port 32000.  Initializing the DB with ./db/init.sql"
podman run --name db -p 32000:3306 -v /tmp/data:/var/lib/mysql -v ./db/init.sql:/docker-entrypoint-initdb.d/init.sql -e MYSQL_ROOT_PASSWORD=root -d docker.io/mysql:8.3.0

echo "You can connect to this db with the command: mysql -u root --host=localhost --port=32000 --protocol=tcp -p"
