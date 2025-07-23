package model

// Campaign, bir yardım kampanyasının bilgilerini tutan veri yapısıdır.
type Campaign struct {
	ID          int
	Title       string
	Description string `gorm:"type:text"` // YENİ EKLENEN ALAN
	Goal        int    // Hedeflenen miktar
	Raised      int    // Toplanan miktar
}
