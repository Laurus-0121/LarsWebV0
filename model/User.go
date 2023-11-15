package model

type User struct {
	ID        int64  `gorm:"primarykey"`
	UserName  string `gorm:"column:user_name;not null" json:"user_name" bind:"required"`
	PassWord  string `grom:"column:pass_word;not null" json:"pass_word" bind:"required"`
	Image     string `gorm:"column:image" json:"image"`
	IsStation string `gorm:"column:is_station;not null" json:"is_station"`
}
