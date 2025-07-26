package model

import "gorm.io/gorm"

// Donation, tamamlanmış bir bağış işlemini temsil eder.
type Donation struct {
    gorm.Model
    DonorName  string `gorm:"type:varchar(100)"`
    DonorEmail string `gorm:"type:varchar(100)"`
    Amount     int    `gorm:"not null"`
    CampaignID uint   // Bu bağışın hangi kampanyaya ait olduğunu belirtir
}