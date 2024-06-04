package services

import (
	"time"

	"github.com/adi-kmt/bitespeed_backend/pkg/entities"
	"github.com/adi-kmt/bitespeed_backend/pkg/repositories"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	repo *repositories.Repository
}

func NewService(repository *repositories.Repository) *Service {
	return &Service{
		repo: repository,
	}
}

func (service *Service) IdentifyUser(ctx *fiber.Ctx, email string, phoneNumber string) (*entities.ContactResponse, *entities.AppError) {
	contacts, err := service.repo.GetAllContactsGivenPhoneOrEmail(ctx, email, phoneNumber)
	if err != nil {
		return nil, err
	}
	if len(contacts) == 0 {
		newContact, err0 := service.repo.InsertContactPrimary(ctx, email, phoneNumber)
		if err0 != nil {
			return nil, err0
		}
		return entities.NewContactResponse(newContact.ID, []string{newContact.Email}, []string{newContact.PhoneNumber}, []int32{}), nil
	} else if len(contacts) == 1 {
		if contacts[0].Email == email && contacts[0].PhoneNumber == phoneNumber {
			return entities.NewContactResponse(contacts[0].ID, []string{contacts[0].Email}, []string{contacts[0].PhoneNumber}, []int32{contacts[0].LinkedID}), nil
		} else {
			newContact, err0 := service.repo.InsertContactSecondary(ctx, email, phoneNumber, contacts[0].ID)
			if err0 != nil {
				return nil, err0
			}
			return entities.NewContactResponse(newContact.ID, []string{newContact.Email}, []string{newContact.PhoneNumber}, []int32{}), nil
		}
	} else {
		var emails []string
		var phoneNumbers []string
		var secondaryIds []int32
		var noOfPrimaries int32 = 0

		type primaryItem struct {
			id        int32
			createdAt time.Time
		}

		var currentPrimary primaryItem

		// First we check if two primaries are present that can have same email or phone number
		for _, contact := range contacts {
			if contact.LinkedID == -1 {
				if noOfPrimaries > 1 {
					if currentPrimary.createdAt.After(contact.CreatedAt) {
						currentPrimary = primaryItem{
							id:        contact.ID,
							createdAt: contact.CreatedAt,
						}
					}
				}
				noOfPrimaries++
			}
		}

		if noOfPrimaries > 0 {
			// Here we convert all the primaries to secondary, except the one created first.
			for _, contact := range contacts {
				if contact.ID != currentPrimary.id {
					newContact, err0 := service.repo.UpdateContact(ctx, currentPrimary.id, email, phoneNumber)
					if err0 != nil {
						return nil, err0
					}
					emails = append(emails, newContact.Email)
					phoneNumbers = append(phoneNumbers, newContact.PhoneNumber)
					secondaryIds = append(secondaryIds, newContact.ID)
				}
			}
		} else {
			for _, contact := range contacts {
				if contact.LinkedID == -1 && (contact.Email == email || contact.PhoneNumber == phoneNumber) {
					//If the item is a primary item, then we append it to the beginning of the list
					currentPrimary = primaryItem{
						id:        contact.ID,
						createdAt: contact.CreatedAt,
					}
					phoneNumbers = append([]string{contact.PhoneNumber}, phoneNumbers...)
					emails = append([]string{contact.Email}, emails...)
				} else {
					// else append to end of the list
					emails = append(emails, contact.Email)
					phoneNumbers = append(phoneNumbers, contact.PhoneNumber)
					secondaryIds = append(secondaryIds, contact.ID)
				}
			}
		}
		return entities.NewContactResponse(currentPrimary.id, emails, phoneNumbers, secondaryIds), nil
	}
}
