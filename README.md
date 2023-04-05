Environment Setup:

Run MariaDB docker container
docker run --detach --name mariadb --env MARIADB_USER=story --env MARIADB_PASSWORD=1q2w3e --env MARIADB_ROOT_PASSWORD=1q2w3e  mariadb:latest --port 3808
