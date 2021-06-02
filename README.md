## Gin Rest Api Example

This project is developed by [Golang](https://golang.org/) with [Gin](https://github.com/gin-gonic/gin) framework to implement simple CRUD with RESTful API.

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

View the swagger page

[Swagger UI](http://localhost:8080/swagger/index.html)

