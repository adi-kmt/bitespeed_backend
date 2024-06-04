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

type primaryItem struct {
	id        int32
	createdAt time.Time
}

func (service *Service) IdentifyUser(ctx *fiber.Ctx, email string, phoneNumber string) (*entities.ContactResponse, *entities.AppError) {
	contacts, err := service.repo.GetAllContactsGivenPhoneOrEmail(ctx, email, phoneNumber)
	if err != nil {
		return nil, err
	}

	var currentPrimary primaryItem
	var emails []string
	var phoneNumbers []string
	var secondaryIds []int32

	if phoneNumber == "" || email == "" {
		checkAndGetNoOfPrimaries(contacts, email, phoneNumber, &currentPrimary)
		return handleSinglePrimary(ctx, service.repo, contacts, currentPrimary, true, email, phoneNumber, &emails, &phoneNumbers, &secondaryIds)
	}

	if len(contacts) == 0 {
		return handleNoContacts(ctx, service.repo, email, phoneNumber)
	}

	isItemAlreadyPresent, noOfPrimaries := checkAndGetNoOfPrimaries(contacts, email, phoneNumber, &currentPrimary)

	if noOfPrimaries > 1 {
		return handleMultiplePrimaries(ctx, service.repo, contacts, currentPrimary, email, phoneNumber)
	}

	return handleSinglePrimary(ctx, service.repo, contacts, currentPrimary, isItemAlreadyPresent, email, phoneNumber, &emails, &phoneNumbers, &secondaryIds)
}

func handleNoContacts(ctx *fiber.Ctx, repo *repositories.Repository, email string, phoneNumber string) (*entities.ContactResponse, *entities.AppError) {
	newContact, err := repo.InsertContactPrimary(ctx, email, phoneNumber)
	if err != nil {
		return nil, err
	}
	return entities.NewContactResponse(newContact.ID, []string{newContact.Email}, []string{newContact.PhoneNumber}, []int32{}), nil
}

func checkAndGetNoOfPrimaries(contacts []entities.ContactDbRecord, email string, phoneNumber string, currentPrimary *primaryItem) (bool, int32) {
	var isItemAlreadyPresent bool = false
	var noOfPrimaries int32 = 0

	for _, contact := range contacts {
		if contact.Email == email && contact.PhoneNumber == phoneNumber {
			isItemAlreadyPresent = true
		}
		if contact.LinkedID == -1 {
			if noOfPrimaries >= 1 {
				if currentPrimary.createdAt.After(contact.CreatedAt) {
					currentPrimary.id = contact.ID
					currentPrimary.createdAt = contact.CreatedAt
				}
			} else {
				currentPrimary.id = contact.ID
				currentPrimary.createdAt = contact.CreatedAt
			}
			noOfPrimaries++
		}
	}
	return isItemAlreadyPresent, noOfPrimaries
}

func handleMultiplePrimaries(ctx *fiber.Ctx, repo *repositories.Repository, contacts []entities.ContactDbRecord, currentPrimary primaryItem, email string, phoneNumber string) (*entities.ContactResponse, *entities.AppError) {
	var emails []string
	var phoneNumbers []string
	var secondaryIds []int32

	for _, contact := range contacts {
		if contact.ID != currentPrimary.id {
			err := repo.UpdateContact(ctx, currentPrimary.id, contact.ID)
			if err != nil {
				return nil, err
			}
		} else {
			secondaryIds = append(secondaryIds, contact.ID)
		}
		emails = append(emails, contact.Email)
		phoneNumbers = append(phoneNumbers, contact.PhoneNumber)
	}

	return entities.NewContactResponse(currentPrimary.id, emails, phoneNumbers, secondaryIds), nil
}

func handleSinglePrimary(ctx *fiber.Ctx, repo *repositories.Repository, contacts []entities.ContactDbRecord, currentPrimary primaryItem, isItemAlreadyPresent bool, email string, phoneNumber string, emails *[]string, phoneNumbers *[]string, secondaryIds *[]int32) (*entities.ContactResponse, *entities.AppError) {
	// If a new item is sent, then it will be inserted as secondary
	if !isItemAlreadyPresent {
		newContact, err := repo.InsertContactSecondary(ctx, email, phoneNumber, currentPrimary.id)
		if err != nil {
			return nil, err
		}
		*phoneNumbers = append(*phoneNumbers, newContact.PhoneNumber)
		*secondaryIds = append(*secondaryIds, newContact.ID)
		*emails = append(*emails, newContact.Email)
	}
	for _, contact := range contacts {
		if contact.LinkedID == -1 {
			*phoneNumbers = append([]string{contact.PhoneNumber}, *phoneNumbers...)
			*emails = append([]string{contact.Email}, *emails...)
		} else {
			*emails = append(*emails, contact.Email)
			*phoneNumbers = append(*phoneNumbers, contact.PhoneNumber)
			*secondaryIds = append(*secondaryIds, contact.ID)
		}
	}

	return entities.NewContactResponse(currentPrimary.id, *emails, *phoneNumbers, *secondaryIds), nil
}
