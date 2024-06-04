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
