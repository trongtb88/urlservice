# urlservice
The service for Shorten URL in Golang

## How to run this service at localhost
1. Start your mysql at your localhost machine successfully
2. Git clone this repo
3. Change file .env which is mapped with your config (DB_USER, DB_PASSWORD, DB_HOST), please note commment DB_HOST at docker, uncomment using 127.0.0.1
4. Create database url_db by yourself
5.  Go to terminal at root of project
```sh
   go get .    
   go run main.go
```

6. If have some logs at console like, server started and workded successfully

```sh

```

7. Go to Postman, import like that, run and get response 200 means work correctly

```
curl --location --request POST 'http://localhost:8023/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url" : "https://gist.github.com/williamn/cfad86ab218101e0c5d7be89226c5c85",
    "shortcode" : "66778b"
}'
```


## How to run this service at docker
1. Run DOCKER DEAMON at your machine successfully
2. Make sure don't have any image mysql is running at port 3306, otherwise you will have 1 error
3. Make sure don't have any image is running at port 8023, otherwise you will have 1 error
2. Git clone this repo
3. Change file .env which is mapped with your config (DB_USER, DB_PASSWORD, DB_HOST), please note uncommment DB_HOST at docker, comment using 127.0.0.1
5.  Go to terminal at root of project
```sh
   docker-compose up --build 
```

6. If have some logs at console like, server started and worked successfully

```sh
url_db_mysql | 2021-06-12T00:58:18.651302Z 0 [Note] Server socket created on IP: '::'.
url_db_mysql | 2021-06-12T00:58:18.652368Z 0 [Warning] Insecure configuration for --pid-file: Location '/var/run/mysqld' in the path is accessible to all OS users. Consider choosing a different directory.
url_db_mysql | 2021-06-12T00:58:18.653396Z 0 [Warning] 'user' entry 'root@url-mysql' ignored in --skip-name-resolve mode.
url_db_mysql | 2021-06-12T00:58:18.661800Z 0 [Note] Event Scheduler: Loaded 0 events
url_db_mysql | 2021-06-12T00:58:18.662171Z 0 [Note] mysqld: ready for connections.
url_db_mysql | Version: '5.7.34'  socket: '/var/run/mysqld/mysqld.sock'  port: 3306  MySQL Community Server (GPL)
shorty_app   | We are connected to the mysql database
shorty_app   | 2021/06/12 00:58:19 /app/api/db/config.go:34
shorty_app   | [1.759ms] [rows:-] SELECT DATABASE()
shorty_app   | 
shorty_app   | 2021/06/12 00:58:19 /app/api/db/config.go:34
shorty_app   | [3.334ms] [rows:-] SELECT count(*) FROM information_schema.tables WHERE table_schema = 'url_db' AND table_name = 'urls' AND table_type = 'BASE TABLE'
shorty_app   | 
shorty_app   | 2021/06/12 00:58:19 /app/api/db/config.go:34
shorty_app   | [8.796ms] [rows:0] CREATE TABLE `urls` (`origin_url` varchar(200),`short_code` varchar(6),`redirect_count` bigint,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`last_seen_at` datetime(3) NULL,PRIMARY KEY (`short_code`))
shorty_app   | 2021/06/12 00:58:19 Listening to port 8023
```

7. Go to Postman, import like that

```
curl --location --request POST 'http://localhost:8023/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url" : "https://gist.github.com/williamn/cfad86ab218101e0c5d7be89226c5c85",
    "shortcode" : "66778b"
}'
```

## How to test this service at localhost

1. Start your mysql at your localhost machine successfully
2. Git clone this repo
3. Change file .env which is mapped with your config (TEST_DB_USER, TEST_DB_PASSWORD, TEST_DB_HOST), please note commment DB_HOST at docker, uncomment using 127.0.0.1
4. Create database test_url_db by yourself
5.  Go to terminal at root of project
```sh
   go test ./tests
```
6. If you see some logs like, mean tests passed.

```sh
 go test ./tests
ok      github.com/trongtb88/urlservice/tests   2.022s

```







