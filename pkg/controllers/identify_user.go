package controllers

import (
	"github.com/adi-kmt/bitespeed_backend/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type identifyUserRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

func IdentifyUser(service *services.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request = new(identifyUserRequest)
		if err := ctx.BodyParser(request); err != nil {
			log.Errorf("error parsing request: %v", err)
			return ctx.Status(fiber.StatusBadRequest).SendString("Error parsing request")
		}
		userData, err1 := service.IdentifyUser(ctx, request.Email, request.PhoneNumber)
		if err1 != nil {
			log.Errorf("error identifying user: %v", err1)
			return ctx.Status(err1.Code).SendString(err1.Error())
		}
		return ctx.Status(200).JSON(userData)
	}
}
