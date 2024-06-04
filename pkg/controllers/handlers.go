package controllers

import (
	"github.com/adi-kmt/bitespeed_backend/pkg/services"
	"github.com/gofiber/fiber/v2"
)

func GetDefaultHandlers(app fiber.Router, service *services.Service) {
	app.Post("/identify", IdentifyUser(service))
}
