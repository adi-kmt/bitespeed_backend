package repositories

import (
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

func (repository *Repository) GetAllContactsGivenPhoneOrEmail(ctx *fiber.Ctx, email string, phoneNumber string) ([]*db.GetContactInfoByEmailORPhoneRow, *entities.AppError) {
	contacts, err := repository.q.GetContactInfoByEmailORPhone(ctx.Context(), db.GetContactInfoByEmailORPhoneParams{
		Column1: email,
		Column2: phoneNumber,
	})
	if err != nil {
		log.Errorf("error getting contacts: %s", err.Error())
		return nil, entities.InternalServerError("Error getting contacts")
	}
	return contacts, nil
}

func (repository *Repository) InsertContact(ctx *fiber.Ctx, email string, phoneNumber string, linkedId int32) (*db.GetContactInfoByEmailORPhoneRow, *entities.AppError) {
	// For primary items
	if linkedId < 0 {
		id, err := repository.q.InsertContactInfo(ctx.Context(), db.InsertContactInfoParams{
			Email:       &email,
			PhoneNumber: &phoneNumber,
		})
		if err != nil {
			log.Errorf("error inserting contacts: %s", err.Error())
			return nil, entities.InternalServerError("Error inserting contacts")
		}
		return &db.GetContactInfoByEmailORPhoneRow{
			ID:          id,
			Email:       &email,
			PhoneNumber: &phoneNumber,
		}, nil
	} else {
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
		return &db.GetContactInfoByEmailORPhoneRow{
			ID:          id,
			Email:       &email,
			PhoneNumber: &phoneNumber,
			LinkedID:    &linkedId,
		}, nil
	}
}

func (repository *Repository) UpdateContactToSecondary(ctx *fiber.Ctx, id int32) (*db.GetContactInfoByEmailORPhoneRow, *entities.AppError) {

}
