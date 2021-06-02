## Gin Rest Api Example

This project is developed by [Golang]([The Go Programming Language (golang.org)](https://golang.org/)) with [Gin]([gin-gonic/gin: Gin is a HTTP web framework written in Go (Golang). It features a Martini-like API with much better performance -- up to 40 times faster. If you need smashing performance, get yourself some Gin. (github.com)](https://github.com/gin-gonic/gin)) framework to implement simple CRUD with RESTful API.

## Prerequisite

Go



## Start from scratch

```shell
go init mod [projectname]
```

Install gin package

```shell
go get -u github.com/gin-gonic/gin
```

Install swagger package

```shell
go get -u github.com/swaggo/swag/cmd/swag
```

Install GORM package

```shell
go get -u gorm.io/gorm
```



To init the Swagger Doc

```shell
swag init
```



Run the app

```shell
go run main.go
```

