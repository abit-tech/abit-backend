package models

type Waitlist struct {
	Email string `gorm:"type:varchar(100);primary_key;not null"`
}

type WaitlistInput struct {
	Email string `json:"email"`
}
