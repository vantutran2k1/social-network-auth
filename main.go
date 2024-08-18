package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/routes"
	"github.com/vantutran2k1/social-network-auth/transaction"
	"github.com/vantutran2k1/social-network-auth/validators"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	validators.RegisterCustomValidators()

	config.InitDB()
	transaction.InitTransactionManager(config.DB)

	router := routes.SetupRouter()
	err = router.Run(":" + os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatal(err)
		return
	}
}
