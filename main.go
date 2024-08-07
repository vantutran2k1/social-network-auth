package main

import (
	"github.com/joho/godotenv"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/routes"
	"github.com/vantutran2k1/social-network-auth/validators"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	validators.RegisterCustomValidators()

	config.InitDB()

	router := routes.SetupRouter()
	err = router.Run(":" + os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatal(err)
		return
	}
}
