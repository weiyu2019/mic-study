package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Mobile   string `gorm:"index:idx_mobile;unique;varchar(11);not null"`
	Password string `gorm:"type:varchar(64);not null"`
	NickName string `gorm:"type:varchar(32)"`
	Gender   string `gorm:"type:varchar(6);default:male"`
	Role     int    `gorm:"type:int;default:1"`
}
