package main

import (
	"fmt" // Package to format strings
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"service_manager/util"
)

func checkAuth(routerContext *gin.Context) {
	fmt.Println("in CheckAuth")
	//check if user is authenticated
	// if not return 401
	// else return 200
	routerContext.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func main() {
	fmt.Println("Service Manager startup initiated...")
	initer, err := util.Initialization("config.yaml")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(initer)
	router := gin.Default()

	// Startup gin Rest Server
	router.GET("/", checkAuth)

	router.Run(":8080")
}
