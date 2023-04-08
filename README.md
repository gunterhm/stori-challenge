# DESCRIPTION

The purpose of the stori-challenge is process incoming files in csv format containing credit and debit transactions for a particular account. The application takes a file from a specific folder, reads all the transactions and send them to a database to be stored. After that, it extracts sata from the database in order to generate a summary report that is sent to the email associated to the account. 

# DESIGN

The solution consists of two main big parts: a Docker container running the Database service which is used to store and persists account data, and a Docker container running the stori-challenge application.

The Database that we use is MariaDB which is a relational SQL database. The database design includes two tables:

- account
  - account_id (PK)
  - name
  - email
- account_transaction
  - account_id (PK)
  - txn_id (PK)
  - amount_credit
  - amount_debit
  - date

The stori-challenge application was created in Go programming language. It was divided into the following components:

- The DB Repository
- The File Processor
- The Report Generator Service
- The Mail Service

The DB Repository Component is in charge of accessing the database for inserting transaction data and for sending SQL queries to the database in order to extract the required data that wil be used to generate the reports. It is using an ORM (Object-Relational Mapper) layer, the chosen ORM is Bun.

The File Processor is the one that finds new files to be processed. Those files need to be in CSV format and  contain the credit and debit transactions that were generated for a particular account. This component reads the content of the CSV file and inserts those transactions into the database. The file processor moves the incoming file to the archive folder in order to avoid the application from processing it more than once. This component also calls the Report Generator and then the Mail Service after a CSV file processing is complete.

The Report Generator Service makes different calls to the Repository in order to get the necessary data from the database for building the Summary Report.

The Mail Service fills in the HTML email template based on the information extracted by the Report  Generation  connects to an SMTP server based on given configuration and sends out the email to the user.

# ENVIRONMENT SETUP

As prerequisite, you need to have Docker installed on your computer. You need to clone this GitHub project on a local folder on your computer.

Before starting to execute the commands below, please make sure you fill in the SMTP password in config.yml file, it was left empty on purpose for security reasons, the password will be provided separately.

We will create and run two containers as explained above, one for the application itself and another one for running the MariaDB database service. Please follow the next steps. 

### Create Docker network to enable communication between containers

- docker network create my_network

### Run MariaDB Docker container

- docker pull mariadb:latest
- docker run --name mariadb_stori --network my_network -e MARIADB_ROOT_PASSWORD=1q2w3e -e MARIADB_USER=stori -e MARIADB_PASSWORD=1q2w3e -e MARIADB_DATABASE=storidb -p 3306:3306 -d mariadb:latest
- docker exec -i mariadb_stori mariadb -ustori -p1q2w3e storidb < resources/initialize.sql

### Run application in Docker container

- docker build -t stori-challenge .
- docker run --name stori-challenge --network my_network -d stori-challenge:latest

