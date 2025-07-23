package model

//import "time"

// FAQ, bir sıkça sorulan soruyu temsil eder.
type FAQ struct {
	ID       uint   `gorm:"primarykey"`
	Question string `gorm:"not null"` // Soru metni
	Answer   string `gorm:"type:text;not null"` // Cevap metni
}
