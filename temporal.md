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

rm development.yaml

ln -s development-mysql8.yaml development.yaml

# 修改数据库的配置
vim development.yaml
```

## create mysql databases && tables
