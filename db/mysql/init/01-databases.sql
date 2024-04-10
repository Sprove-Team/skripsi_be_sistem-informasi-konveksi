# create databases
CREATE DATABASE IF NOT EXISTS `db_khris_production`;
CREATE DATABASE IF NOT EXISTS `db_khris_production_test`;

# create root user and grant rights
CREATE USER 'root'@'localhost' IDENTIFIED BY 'root';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%';