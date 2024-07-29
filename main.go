package main

import (
	"github.com/joho/godotenv"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/routes"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	appPort := os.Getenv("APP_PORT")

	router := routes.SetupRouter()

	err = router.Run(":" + appPort)
	if err != nil {
		log.Fatal(err)
		return
	}
}
