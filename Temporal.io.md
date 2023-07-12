# Temporal.io

## Temporal Install

### Download
[Temporal Server](https://github.com/temporalio/temporal/releases)

```shell
tar -xvf temporal_1.21.0_darwin_amd64.tar.gz

cd /path/to/temporal_1.21.0_darwin_amd64
```

## MySQL8.0 Docker Image
```shell
# https://hub.docker.com/_/mysql
# run 
docker run --name temporal -v /path/to/temporal-data-dir:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=temporal -p 3306:3306 -d mysql:8.0

# bash
docker exec -it temporal bash

# mysql
mysql -h 127.0.0.1  -uroot -ptemporal
```

## Temporal Config
```shell
cd /path/to/temporal_1.21.0_darwin_amd64/config

# if exists, `rm` it
rm development.yaml

#ln
ln -s development-mysql8.yaml development.yaml

# update database info
vim development.yaml
```

## Create MySQL8 Databases && Tables
```shell
cd /path/to/temporal_1.21.0_darwin_amd64/

git clone -bv1.21.0 --depth=1 https://github.com/temporalio/temporal.git

# https://github.com/temporalio/docker-builds/blob/main/docker/auto-setup.sh
# setup_mysql_schema()

# create database temporal && tables
./temporal-sql-tool --endpoint 127.0.0.1 --port 3306 --user root --password temporal --database temporal --plugin mysql create

./temporal-sql-tool \
--endpoint 127.0.0.1 --port 3306 --user root --password temporal  --database temporal --plugin mysql \
setup-schema --version 0.0

 ./temporal-sql-tool \
--endpoint 127.0.0.1 --port 3306 --user root --password temporal  --database temporal --plugin mysql \
update-schema --schema-dir ./temporal/schema/mysql/v8/temporal/versioned

# create database temporal_visibility && tables
./temporal-sql-tool --endpoint 127.0.0.1 --port 3306 --user root --password temporal --database temporal_visibility --plugin mysql create

./temporal-sql-tool --endpoint 127.0.0.1 --port 3306 --user root --password temporal --database temporal_visibility --plugin mysql setup-schema -v 0.0

./temporal-sql-tool \
--endpoint 127.0.0.1 --port 3306 --user root --password temporal  --database temporal_visibility --plugin mysql \
update-schema --schema-dir ./temporal/schema/mysql/v8/visibility/versioned
```


## start server
