package entities

type ContactResponse struct {
	Contact Contact `json:"contact"`
}

type Contact struct {
	PrimaryContactId        int32    `json:"primaryContactId"`
	Emails                  []string `json:"emails"`
	PhoneNumbers            []string `json:"phoneNumbers"`
	SecondaryContactNumbers []int32  `json:"secondaryContactNumbers"`
}

func NewContactResponse(primaryContactId int32, emails []string, phoneNumbers []string, secondaryContactNumbers []int32) *ContactResponse {
	return &ContactResponse{
		Contact{
			PrimaryContactId:        primaryContactId,
			Emails:                  emails,
			PhoneNumbers:            phoneNumbers,
			SecondaryContactNumbers: secondaryContactNumbers,
		},
	}
}
