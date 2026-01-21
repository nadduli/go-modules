package models

type User struct {
	BaseModel

	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}
