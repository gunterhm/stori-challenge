Environment Setup:


Run MariaDB docker container

docker pull mariadb:latest

docker run –name mariadb_stori -e MYSQL_ROOT_PASSWORD=1q2w3e -p 3306:3306 -d docker.io/library/mariadb:10.3
