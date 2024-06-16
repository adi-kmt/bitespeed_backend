package services

import (
	"github.com/adi-kmt/bitespeed_backend/pkg/entities"
	"github.com/adi-kmt/bitespeed_backend/pkg/repositories"
	"github.com/adi-kmt/bitespeed_backend/pkg/utils"
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
	primaryContacts, err := service.repo.GetPrimaryContactsGivenEmailOrPhone(ctx, email, phoneNumber)
	if err != nil {
		return nil, err
	}

	var secondaryContacts []entities.ContactDbRecord
	for _, contact := range primaryContacts {
		contacts, err0 := service.repo.GetSecondaryContacts(ctx, contact.ID)
		if err0 != nil {
			return nil, err0
		}
		secondaryContacts = append(secondaryContacts, contacts...)
	}

	if phoneNumber == "" || email == "" {
		return handleSinglePrimary(ctx, service.repo, primaryContacts[0], secondaryContacts, email, phoneNumber, true)
	}

	if len(primaryContacts) == 0 {
		return handleNoContacts(ctx, service.repo, email, phoneNumber)
	} else if len(primaryContacts) == 1 {
		return handleSinglePrimary(ctx, service.repo, primaryContacts[0], secondaryContacts, email, phoneNumber, false)
	} else {
		return handleMultiplePrimaries(ctx, service.repo, primaryContacts, secondaryContacts, email, phoneNumber)
	}
}

func handleNoContacts(ctx *fiber.Ctx, repo *repositories.Repository, email string, phoneNumber string) (*entities.ContactResponse, *entities.AppError) {
	newContact, err := repo.InsertContactPrimary(ctx, email, phoneNumber)
	if err != nil {
		return nil, err
	}
	return entities.NewContactResponse(newContact.ID, []string{newContact.Email}, []string{newContact.PhoneNumber}, []int32{}), nil
}

func handleMultiplePrimaries(ctx *fiber.Ctx, repo *repositories.Repository, primaries []entities.ContactDbRecord, secondaries []entities.ContactDbRecord, newEmail, newPhoneNo string) (*entities.ContactResponse, *entities.AppError) {
	var actualPrimary entities.ContactDbRecord = primaries[0]
	var newSecondary entities.ContactDbRecord
	var isEmailPresent bool = false
	var isPhoneNoPresentInAnother bool = false

	for _, primary := range primaries {
		if primary.Email == newEmail {
			isEmailPresent = true
		}
		if primary.PhoneNumber == newPhoneNo {
			isPhoneNoPresentInAnother = true
		}
	}

	if actualPrimary.CreatedAt.After(primaries[1].CreatedAt) {
		actualPrimary = primaries[1]
		secondaries = append(secondaries, primaries[0])
		newSecondary = primaries[0]
	} else {
		secondaries = append(secondaries, primaries[1])
		newSecondary = primaries[1]
	}

	err := repo.UpdateContact(ctx, actualPrimary.ID, newSecondary.ID)
	if err != nil {
		return nil, err
	}

	return handleSinglePrimary(ctx, repo, actualPrimary, secondaries, newEmail, newPhoneNo, isEmailPresent && isPhoneNoPresentInAnother)
}

func handleSinglePrimary(ctx *fiber.Ctx, repo *repositories.Repository, primary entities.ContactDbRecord, secondaries []entities.ContactDbRecord, newEmail, newPhoneNo string, checkPresent bool) (*entities.ContactResponse, *entities.AppError) {
	contactEntity := entities.NewContactResponse(primary.ID, []string{primary.Email}, []string{primary.PhoneNumber}, []int32{})
	if primary.Email == newEmail && primary.PhoneNumber == newPhoneNo && !checkPresent {
		checkPresent = true
	}

	for _, secondary := range secondaries {
		contactEntity.Contact.Emails = append(contactEntity.Contact.Emails, secondary.Email)
		contactEntity.Contact.PhoneNumbers = append(contactEntity.Contact.PhoneNumbers, secondary.PhoneNumber)
		contactEntity.Contact.SecondaryContactNumbers = append(contactEntity.Contact.SecondaryContactNumbers, secondary.ID)
		if secondary.Email == newEmail && secondary.PhoneNumber == newPhoneNo && !checkPresent {
			checkPresent = true
		}
	}

	if !checkPresent {
		contact, err := repo.InsertContactSecondary(ctx, newEmail, newPhoneNo, primary.ID)
		if err != nil {
			return nil, err
		}
		contactEntity.Contact.SecondaryContactNumbers = append(contactEntity.Contact.SecondaryContactNumbers, contact.ID)
		contactEntity.Contact.Emails = append(contactEntity.Contact.Emails, contact.Email)
		contactEntity.Contact.PhoneNumbers = append(contactEntity.Contact.PhoneNumbers, contact.PhoneNumber)
	}
	contactEntity.Contact.Emails = utils.UniqueSliceElements(contactEntity.Contact.Emails)
	contactEntity.Contact.PhoneNumbers = utils.UniqueSliceElements(contactEntity.Contact.PhoneNumbers)
	contactEntity.Contact.SecondaryContactNumbers = utils.UniqueSliceElements(contactEntity.Contact.SecondaryContactNumbers)
	return contactEntity, nil
}
