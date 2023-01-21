# CHECK24 in Go

## Build & Deployment
```shell
$env:GOARCH = "arm"
$env:GOARM = "6"
$env:GOOS = "linux"
# $env:GOARCH = "amd64"
# $env:GOOS = "windows"
go env

go build check24-in-go/cmd/server
go build check24-in-go/cmd/dbinit


rsync -az -e ssh check24-in-go [user]@[ip]:[directory]
rsync -az -e ssh public [user]@[ip]:[directory]/public

rsync -az -e ssh server test@172.30.124.56:/home/test
rsync -az -e ssh public test@172.30.124.56:/home/test
```

## Database Setup
```shell
sudo apt-get update
sudo apt-get upgrade -y
sudo apt-get install -y mariadb-server
sudo mysql_secure_installation
sudo mysql -u root -p
```

The length of the column is printed on startup of the dbinit command.

```sql
CREATE DATABASE products;
CREATE USER 'go'@'localhost' IDENTIFIED BY 'password';
CREATE USER 'go'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON products.* TO 'go'@'localhost';
GRANT ALL PRIVILEGES ON products.* TO 'go'@'%';
FLUSH PRIVILEGES;

CREATE TABLE products(NAME VARCHAR(220), PRICE INT, IMAGE VARCHAR(479));


CREATE DATABASE logging;
GRANT ALL PRIVILEGES ON logging.* TO 'go'@'localhost';
GRANT ALL PRIVILEGES ON logging.* TO 'go'@'%';

CREATE TABLE logging(URL TEXT, USERAGENT TEXT, TIME DATETIME);
```
