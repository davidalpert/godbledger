## Setting up MySQL
On Ubuntu 18.04

### Install
This will allow you to login to the root with `sudo mysql`

```
sudo apt-get install mysql-server
sudo service mysql start
mysql_secure_installation
sudo mysql
```

To stop the service use
```
sudo service mysql stop
```

### Create a database for GoDBLedger
While logged in as root
```
CREATE DATABASE ledger;
```

### Create User for GoDBLedger
Login to the mysql prompt
```
CREATE USER 'godbledger'@'localhost' IDENTIFIED BY 'password';
```
This will create the following user:
**User:** godbledger
**Pass:** password

Grant the GoDBLedger user rights to databases & logins etc
```
GRANT ALL PRIVILEGES ON * . * TO 'godbledger'@'localhost';
FLUSH PRIVILEGES;
```
Please note that in this example we are granting godbledger full root access to everything in our database. While this is helpful for explaining some concepts, it may be impractical for most use cases and could put your database’s security at high risk.

### Login to the ledger database with your GoDBLedger user
```
mysql -u godbledger -ppassword ledger
```

### Delete the User (Optional)
From within MySQL prompt
```
DROP USER ‘godbledger’@‘localhost’;
```

## Allowing External Connections to MySQL

Setup the user to allow external IP addresses to access
```
GRANT ALL PRIVILEGES ON *.* TO 'USERNAME'@'IP' IDENTIFIED BY 'PASSWORD' with grant option;
FLUSH PRIVILEGES;
```
So for our godbledger user it becomes
```
GRANT ALL PRIVILEGES ON *.* TO 'godbledger'@'%' IDENTIFIED BY 'password' with grant option;
```

Then you need to allow external connections to MySQL as a whole. Find the configuration file then find the following line and comment it out in your `my.cnf file`, which usually lives on `/etc/mysql/my.cnf` on Unix/OSX systems. In some cases the location for the file is `/etc/mysql/mysql.conf.d/mysqld.cnf`).

Can also be found using
```
sudo find /etc -iname 'mysql*.cnf'
```

Change line
```
bind-address = 127.0.0.1
```

to
```
#bind-address = 127.0.0.1
```
