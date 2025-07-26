package model

import "gorm.io/gorm"

// Message, iletişim formundan gelen mesajları temsil eder.
type Message struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	Email   string `gorm:"type:varchar(100);not null"`
	Subject string `gorm:"type:varchar(255);not null"`
	Message string `gorm:"type:text;not null"`
}