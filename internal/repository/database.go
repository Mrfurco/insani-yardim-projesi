package repository

import (
	"fmt"
	"log"

	"github.com/MrFurco/insani-yardim-projesi/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase, PostgreSQL veritabanına bağlanır.
func ConnectDatabase() {
	// docker-compose.yml dosyasında belirlediğimiz bilgilerle bağlantı cümlesi oluşturuyoruz.
	dsn := "host=localhost user=yardim_user password=strong_password_123 dbname=insani_yardim_db port=5432 sslmode=disable"
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı!", err)
	}

	fmt.Println("Veritabanına başarıyla bağlanıldı!")

	// Otomatik olarak 'campaigns' tablosunu oluşturur (eğer yoksa).
	// Bu, model.Campaign yapısına bakarak tabloyu oluşturur.
	err = database.AutoMigrate(&model.Campaign{}, &model.Post{}, &model.FAQ{}, &model.Message{}, &model.Donation{})
	if err != nil {
		log.Fatal("Tablo oluşturulamadı!", err)
	}

	

	fmt.Println("Veritabanı tablosu başarıyla oluşturuldu/güncellendi.")
	DB = database
}