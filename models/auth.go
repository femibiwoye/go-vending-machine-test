package models

type User struct {
	ID         uint   `gorm:"primaryKey" json:"id,omitempty"`
	FullName   string `json:"full_name,omitempty"`
	UserName   string `json:"user_name,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `gorm:"index" json:"phone,omitempty"`
	Password   string `json:"password,omitempty"`
	CreatedAt  int64  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt  int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	IsVerified bool   `json:"is_verified,omitempty"`
	Role       string `json:"role,omitempty"`
	Deposit    int    `json:"deposit"`
}

type Session struct {
	ID     uint   `gorm:"primaryKey" json:"id,omitempty"`
	UserID uint   `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

type AuthCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type UserUpdate struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone" `
	Role     string `json:"role" `
}
