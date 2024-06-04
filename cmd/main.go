package main

import (
	"fmt"
	"os"

	"github.com/adi-kmt/bitespeed_backend/pkg/controllers"
	"github.com/adi-kmt/bitespeed_backend/pkg/injection"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Main File %s", err.Error())
	}

	service := injection.InjectDependencies()

	controllers.GetDefaultHandlers(app, service)
	port, isErr := os.LookupEnv("API_PORT")
	if !isErr {
		port = "8080"
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
