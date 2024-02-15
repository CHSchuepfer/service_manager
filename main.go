package main

import (
	"./util"
	"fmt" // Package to format strings
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
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
	fmt.Println("Service Manin CheckAuthager startup initiated...")
	util.InitRuntime()
	router := gin.Default()

	// Startup gin Rest Server
	router.GET("/", checkAuth)

	router.Run(":8080")
}
