package entities

import (
	"time"

	db "github.com/adi-kmt/bitespeed_backend/db/sqlc"
)

type ContactDbRecord struct {
	ID          int32
	Email       string
	PhoneNumber string
	LinkedID    int32
	CreatedAt   time.Time
}

func NewContactDbRecord(id int32, email string, phoneNumber string, linkedID int32, createdAt time.Time) *ContactDbRecord {
	return &ContactDbRecord{
		ID:          id,
		Email:       email,
		PhoneNumber: phoneNumber,
		LinkedID:    linkedID,
		CreatedAt:   createdAt,
	}
}

func CreateRecordFromGetContact(dbRow []*db.GetContactInfoByEmailORPhoneRow) []ContactDbRecord {
	var records []ContactDbRecord
	for _, contact := range dbRow {
		var linkedIDString int32
		if contact.LinkedID != nil {
			linkedIDString = *contact.LinkedID
		} else {
			linkedIDString = -1
		}
		records = append(records, ContactDbRecord{
			ID:          contact.ID,
			Email:       *contact.Email,
			PhoneNumber: *contact.PhoneNumber,
			LinkedID:    linkedIDString,
			CreatedAt:   contact.CreatedAt.Time,
		})
	}
	return records
}
