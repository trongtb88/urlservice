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
   go main run
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
   cd tests
   go test ./tests
```
6. If you see some logs like, mean tests passed.

```sh
 go test ./tests
ok      github.com/trongtb88/urlservice/tests   2.022s


```







