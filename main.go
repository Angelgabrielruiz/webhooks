package main

import (
	"github_wb/infrastructure"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	router := gin.Default()

	infrastructure.Routes(router)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	//a ver cxon web 3
	router.Run(":" + port)

}
