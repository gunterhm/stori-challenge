Environment Setup:


Run MariaDB docker container

docker pull mariadb:latest

docker run --name mariadb_stori -e MARIADB_ROOT_PASSWORD=1q2w3e -e MARIADB_USER=stori -e MARIADB_PASSWORD=1q2w3e -e MARIADB_DATABASE=storidb -p 3306:3306 -d mariadb:latest

docker exec -i mariadb_stori mariadb -ustori -p1q2w3e storidb < resources/initialize.sql

