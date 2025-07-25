package model

import "time"

// Post, bir haber veya blog yazısının bilgilerini tutar.
type Post struct {
	ID        uint      `gorm:"primarykey"`
	Title     string    // Yazının başlığı
	Summary   string    `gorm:"type:text"`
	Content   string    `gorm:"type:text"` // Yazının içeriği (uzun metin)
	CreatedAt time.Time // GORM tarafından otomatik olarak oluşturulma zamanı eklenecek
	UpdatedAt time.Time // GORM tarafından otomatik olarak güncellenme zamanı eklenecek
	ImageURL  string    `gorm:"type:varchar(255)"`
}
