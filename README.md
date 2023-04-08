ENVIRONMENT SETUP


docker network create my_network

RUN MARIADB DOCKER CONTAINER

docker pull mariadb:latest

docker run --name mariadb_stori --network my_network -e MARIADB_ROOT_PASSWORD=1q2w3e -e MARIADB_USER=stori -e MARIADB_PASSWORD=1q2w3e -e MARIADB_DATABASE=storidb -p 3306:3306 -d mariadb:latest

docker exec -i mariadb_stori mariadb -ustori -p1q2w3e storidb < resources/initialize.sql



RUN APPLICATION IN DOCKER CONTAINER

docker build -t stori-challenge .

docker run --name stori-challenge --network my_network -d stori-challenge:latest

