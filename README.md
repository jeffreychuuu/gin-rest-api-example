# Gin Rest Api Example

This project is developed by [Golang](https://golang.org/) with [Gin](https://github.com/gin-gonic/gin) framework to implement simple CRUD with RESTful API.

- [Gin Web Framework](#https://github.com/gin-gonic/gin) - A martini-like API with performance that is up to 40 times faster thanks to httprouter
- [GORM](#https://gorm.io/index.html) - The fantastic ORM library for Golang
- [Go Redis](#https://github.com/go-redis/redis) - Supports 2 last Go versions and requires support for Go modules
- [gRPC](#https://grpc.io/docs/languages/go/quickstart/) - A modern open source high performance Remote Procedure Call (RPC) framework that can run in any environment

## Table of Contents

- [Prerequisite](#Prerequisite)
- [Start from scratch](#Start-from-scratch)
- [Get Started](#Get-Started)
- [Integrate with Gin Framework](#Integrate-with-Gin-Framework)
- [Integrate with GoDotEnv](#Integrate-with-GoDotEnv)
- [Integrate with GORM](#Integrate-with-GORM)
- [Integrate with Go Redis](#Integrate-with-Go-Redis)
- [Integrate with gRPC](#Integrate-with-gRPC)



## Prerequisite

Go



## Start from scratch

```shell
go mod init [projectname]
```



## Get Started

Run the app

```shell
go run main.go
```

View the swagger page

[Swagger UI](http://localhost:8080/swagger/index.html)



## Integrate with Gin Framework

Install gin package

```shell
go get -u github.com/gin-gonic/gin
```

Define the router for each rest api in main.go

```go
package main

import (
	"gin-rest-api-example/controllers"
	_ "gin-rest-api-example/docs"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Routes
	// Book Router
	bookRouter := r.Group("")
	{
		bookRouter.GET("/books", controllers.FindBooks)
		bookRouter.GET("/books/:id", controllers.FindBook)
		bookRouter.POST("/books", controllers.CreateBook)
		bookRouter.PATCH("/books/:id", controllers.UpdateBook)
		bookRouter.DELETE("/books/:id", controllers.DeleteBook)
	}

	// Run the server
	r.Run()
}
```



## Integrate with GoDotEnv

Install go dot env package as library

```sh
go get github.com/joho/godotenv
```

or if you want to use it as a bin command

```sh
go get github.com/joho/godotenv/cmd/godotenv
```

Define the .env file

```env
POSTGRES_HOST=localhost
```

Create an autoload method to load .env file once application started

```go
package autoload

/*
   You can just read the .env file on import just by doing
       import _ "github.com/joho/godotenv/autoload"
   And bob's your mother's brother
*/

import "github.com/joho/godotenv"

func init() {
	godotenv.Load()
}
```

Access the environment variables everywhere

```go
os.Getenv("POSTGRES_HOST")
```



## Integrate with GORM

Install GORM package

```shell
go get -u gorm.io/gorm
```

Create the db connection client

```go
package client

import (
	"fmt"
	"gin-rest-api-example/models"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

var DB *gorm.DB

func ConnectDatabase() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
	)
	database, err := gorm.Open("postgres", connStr)

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&models.Book{})

	DB = database
}
```

Make the connection in main.go when the application started

```go
package main

import (
  "gin-rest-api-example/client"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Connect to database
	client.ConnectDatabase()

	// Run the server
	r.Run()
}
```

Call the db client for CRUD action in services

```go
package services

import (
	"gin-rest-api-example/client"
	"gin-rest-api-example/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FindBooks(c *gin.Context) {
	var books []models.Book
  // use the export DB
	client.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}
```



## Integrate with Go Redis

Install Go Redis package

```sh
go get github.com/go-redis/redis/v8
```

Create the Redis connection client

```go
package client

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

var Redis *redis.Client
var ctx = context.Background()

func ConnectRedis() *redis.Client {
	address := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Connection fail in Redis：", pong, err)
		panic(err)
	}
	fmt.Println("Connection success in Redis：", pong)
	Redis = client
	return client
}
```

Make the connection in main.go when the application started

```go
package main

import (
  "gin-rest-api-example/client"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Connect to redis
	client.ConnectRedis()

	// Run the server
	r.Run()
}
```

Call the Redis client in services

```go
package services

import (
	"gin-rest-api-example/client"
	"gin-rest-api-example/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FindBook(c *gin.Context) {
	var book models.Book
	// Try getting from Redis
	bookJson, _ := client.Redis.HGet(c, "Book", c.Param("id")).Result()
	json.Unmarshal([]byte(bookJson), &book)

	// Get model if exist
	if bookJson == "" {
		// Get Book from DB
		if err := client.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}
		// Cache Book in Redis
		bookJson, err := json.Marshal(book)
		err = client.Redis.HSet(c, "Book", strconv.FormatUint(uint64(book.ID), 10), string(bookJson)).Err()
		ttl, err := strconv.ParseInt(os.Getenv("REDIS_TTL"), 10, 64)
		// Set Redis Expire Time
		client.Redis.Expire(c, "Book", time.Duration(ttl)*time.Second)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Redis Insertion Success!")
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}
```





## Integrate with gRPC

Install swagger package

```shell
go get -u github.com/swaggo/swag/cmd/swag
```

Define the Openapi document in main.go

```go
package main

import (
	_ "gin-rest-api-example/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gin Rest Api Example Swagger
// @version 1.0
// @description Gin Rest Api Example Swagger

// @contact.name Jeffrey Chu
// @contact.email jeffreychu888hk@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
func main() {
	r := gin.Default()

	// Swagger
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Run the server
	r.Run()
}
```

Define each api method in controller

```go
package controllers

import (
	"gin-rest-api-example/services"
	"github.com/gin-gonic/gin"
)

// @Tags Book
// @Summary Find books
// @Success 200 {object} models.Result Successful Return Value
// @Router /books [get]
func FindBooks(c *gin.Context) {
	services.FindBooks(c)
}
```



To init the Swagger Doc

```shell
swag init
```



## Integrate with gRPC

Install Go protocol buffers plugin

````sh
go get github.com/golang/protobuf/protoc-gen-go
````

Install Golang grpc package

```sh
go get google.golang.org/grpc
```

Install Go Micro

```sh
go get github.com/micro/micro/v3
```



Define a proto file

```protobuf
syntax = "proto3";

package book;
option go_package = "./";

service BookService{ 
  rpc CreateBook (CreateBookInput) returns (Result) {}
}

message CreateBookInput {
  string title = 1;
  string author = 2;
  bool isEnable = 3;
}

message Result {
  int32 code = 1;
  string message = 2;
  string data = 3;
}
```

Using protoc cli to generate a .pb.go file which may contain the services, functions, requests and response that defined in the protofile

```sh
protoc --go_out=plugins=grpc:. *.proto
```

Import the pb file that generated before 

*Ensure that using the full path if the pb file is not from github*

```go
//Full Path ~/gin-rest-api-example/proto/book.pb
import pb "gin-rest-api-example/proto"
```

### Create gRPC Server

Define the gRPC Server

```go
package grpc

import (
	"context"
	"gin-rest-api-example/models"
	pb "gin-rest-api-example/proto"
	"gin-rest-api-example/services"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func StartGrpcServer() {
	// Create gRPC Server
	const host = "localhost"
	const port = "5000"
	lis, err := net.Listen("tcp", host + ":" + port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
  // Remeber to use goroutine to run gRPC as microservice
	go func() {
		s := grpc.NewServer()
		log.Printf("gRPC server is running in port: %s.", port)

		pb.RegisterBookServiceServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (s server) CreateBook(c context.Context, input *pb.CreateBookInput) (*pb.Result, error) {
	request := models.CreateBookInput{Title: input.Title, Author: input.Author, IsEnable: input.IsEnable}
	response := services.CreateBook(request)

	return &pb.Result{Code: int32(response.Code), Message: response.Message,Data: response.Data.(string)}, nil
}
```

Start the gRPC server in main.go

```go
package main

import (
	"gin-rest-api-example/grpc"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()

	// Start grpc server
	grpc.StartGrpcServer()

	// Run the gin server
	r.Run()
}
```



### Create gRPC Client

Connect to the gRPC Client

```go
package main

import (
	"context"
	"gin-rest-api-example/models"
	pb "gin-rest-api-example/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBookServiceClient(conn)

	input := models.CreateBookInput{
		Title: "Test", Author: "Jeffrey", IsEnable: true,
	}
	createBook(c, input)
}

func createBook(c pb.BookServiceClient, input models.CreateBookInput) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.CreateBook(ctx, &pb.CreateBookInput{Author: input.Author, Title: input.Title, IsEnable: input.IsEnable})
	if err != nil {
		log.Fatalf("Could not createBook: %v", err)
	}
	log.Printf("gRPC response: %s", res.Data)
}
```

### Try to run simulate a gRPC Server and Client

```sh
// start the gRPC Server
go run main.go

// start the gRPC Client
go run client.go
```

