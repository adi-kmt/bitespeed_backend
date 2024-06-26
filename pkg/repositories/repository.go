package repositories

import (
	"time"

	db "github.com/adi-kmt/bitespeed_backend/db/sqlc"
	"github.com/adi-kmt/bitespeed_backend/pkg/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q    *db.Queries
	conn *pgxpool.Pool
}

func NewRepository(q *db.Queries, conn *pgxpool.Pool) *Repository {
	return &Repository{
		q:    q,
		conn: conn,
	}
}

func (repository *Repository) GetPrimaryContactsGivenEmailOrPhone(ctx *fiber.Ctx, email string, phoneNumber string) ([]entities.ContactDbRecord, *entities.AppError) {
	contacts, err := repository.q.GetPrimaryContactInfoByEmailORPhone(ctx.Context(), db.GetPrimaryContactInfoByEmailORPhoneParams{
		Email:       &email,
		PhoneNumber: &phoneNumber,
	})
	if err != nil {
		log.Errorf("error getting contacts: %s", err.Error())
		return nil, entities.InternalServerError("Error getting contacts")
	}
	return entities.CreateRecordFromGetContact(contacts), nil
}

func (repository *Repository) GetSecondaryContacts(ctx *fiber.Ctx, primaryId int32) ([]entities.ContactDbRecord, *entities.AppError) {
	contacts, err := repository.q.GetSecondaryContactInfo(ctx.Context(), &primaryId)
	if err != nil {
		log.Errorf("error getting contacts: %s", err.Error())
		return nil, entities.InternalServerError("Error getting contacts")
	}
	return entities.CreateRecordFromGetContactSecondary(contacts), nil
}

func (repository *Repository) InsertContactPrimary(ctx *fiber.Ctx, email string, phoneNumber string) (*entities.ContactDbRecord, *entities.AppError) {
	id, err := repository.q.InsertContactInfo(ctx.Context(), db.InsertContactInfoParams{
		Email:          &email,
		PhoneNumber:    &phoneNumber,
		LinkedID:       nil,
		LinkPrecedence: db.LinkPrecedenceEnumPrimary,
	})
	if err != nil {
		log.Errorf("error inserting contacts: %s", err.Error())
		return nil, entities.InternalServerError("Error inserting contacts")
	}
	return entities.NewContactDbRecord(id, email, phoneNumber, -1, time.Now()), nil
}

func (repository *Repository) InsertContactSecondary(ctx *fiber.Ctx, email string, phoneNumber string, linkedId int32) (*entities.ContactDbRecord, *entities.AppError) {
	id, err := repository.q.InsertContactInfo(ctx.Context(), db.InsertContactInfoParams{
		Email:          &email,
		PhoneNumber:    &phoneNumber,
		LinkedID:       &linkedId,
		LinkPrecedence: db.LinkPrecedenceEnumSecondary,
	})
	if err != nil {
		log.Errorf("error inserting contacts: %s", err.Error())
		return nil, entities.InternalServerError("Error inserting contacts")
	}
	return entities.NewContactDbRecord(id, email, phoneNumber, linkedId, time.Now()), nil
}

func (repository *Repository) UpdateContact(ctx *fiber.Ctx, linkedId, id int32) *entities.AppError {
	err := repository.q.UpdateContactToSecondary(ctx.Context(), db.UpdateContactToSecondaryParams{
		LinkedID: &linkedId,
		ID:       id,
	})
	if err != nil {
		log.Errorf("error updating contacts: %s", err.Error())
		return entities.InternalServerError("Error updating contacts")
	}
	return nil
}
