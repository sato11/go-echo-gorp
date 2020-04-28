# go-echo-gorp
A sample application to try out some tools, namely:
- Web framework: [github.com/labstack/echo](https://github.com/labstack/echo), and
- ORMapper: [github.com/go-gorp/gorp](https://github.com/go-gorp/gorp)

## Run
```
## Start up containers
$ docker-compose build
$ docker-compose up -d

## Call APIs from your host machine
# GET
$ curl localhost:8080/api/comments
# or POST
$ curl -X POST -H 'Content-Type: application/json' -d '{"text":"Hello World!"}' localhost:8080/api/comments
```

## Test
```
## Use docker-compose.test.yml
$ docker-compose -f docker-compose.test.yml up -d
# or simply
$ docker-compose -f docker-compose.test.yml run api make test
```
