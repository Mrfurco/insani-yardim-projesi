package main

import (
	"log"

	"github.com/MrFurco/insani-yardim-projesi/internal/handler"
	"github.com/MrFurco/insani-yardim-projesi/internal/model"
	"github.com/MrFurco/insani-yardim-projesi/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	repository.ConnectDatabase()
	seedDatabase(repository.DB)

	router := gin.Default()
	//router.Static("/static", "./static")

	// Şablonları en basit ve güvenilir şekilde yüklüyoruz.
	router.LoadHTMLGlob("templates/*")

	// --- Public Rotalar ---
	router.GET("/", handler.ListCampaigns)
	router.GET("/kampanya/:id", handler.ShowCampaign)

	// --- Admin Grubu ve Rotaları ---
	adminGroup := router.Group("/admin")
	adminGroup.Use(gin.BasicAuth(gin.Accounts{
		"admin": "12345",
	}))
	// Kampanya Yönetimi
	adminGroup.GET("/dashboard", handler.ShowAdminDashboard)
	adminGroup.GET("/kampanyalar/yeni", handler.ShowNewCampaignForm)
	adminGroup.POST("/kampanyalar", handler.CreateCampaign)
	adminGroup.GET("/kampanyalar/sil/:id", handler.DeleteCampaign)
	adminGroup.GET("/kampanyalar/duzenle/:id", handler.ShowEditCampaignForm)
	adminGroup.POST("/kampanyalar/duzenle/:id", handler.UpdateCampaign)
	// Haber Yönetimi
	adminGroup.GET("/haberler", handler.ListPosts)
	adminGroup.GET("/haberler/yeni", handler.ShowNewPostForm)
	adminGroup.POST("/haberler", handler.CreatePost)

	router.Run()
}

// seedDatabase, veritabanı boşsa başlangıç verilerini ekler.
func seedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&model.Campaign{}).Count(&count)

	if count == 0 {
		log.Println("Veritabanı boş, gerçek başlangıç verileri ekleniyor...")

		campaigns := []model.Campaign{
			{
				Title:       "Filistin İçin Acil Yardım",
				Description: `Filistin'in 7 Ekim'de işgalci İsrail'e karşı başlattığı "Aksa Tufanı" özgürlük mücadelesi devam ediyor. İsrail su ve elektrik hatlarını kestiği Filistin'de şu ana kadar bin tondan fazla bomba kullandığını açıkladı. Tüm bu saldırılar sonucunda hastane, cami ve okulların da bulunduğu birçok yerleşim yeri ağır hasar aldı. Gazze'ye düzenlenen bu saldırıların ardından yapılan son açıklamalara göre 17 binden den fazlası çocuk olmak üzere 50 binden fazla Filistinli vefat ederken 113 binden fazla kişinin yaralandığı belirtildi.`,
				Goal:        500000,
				Raised:      95000,
			},
			{
				Title:       "Bir Burs Bir Gelecek",
				Description: `Geleceğimizi güvence altına almanın en güzel yolu eğitimdir. Hayrat Yardım olarak, ülkemizdeki öğrencilerin eğitim hayatlarına devam edebilmeleri için "Bir Burs Bir Gelecek" kampanyasını başlattık. Bağışlarınızla bir öğrencinin eğitim masraflarına ortak olabilir, onun geleceğine ışık tutabilirsiniz.`,
				Goal:        100000,
				Raised:      85000,
			},
			{
				Title:       "Bu Kış Kimse Üşümesin",
				Description: `Yurt içinde ve yurt dışında, kışın ağır koşullarında yaşam mücadelesi veren birçok aile, soğuktan korunmakta zorluk çekiyor. Özellikle çocuklar, koruyucu kıyafet eksikliği nedeniyle zarar görüyor. Dağıttığımız kışlık kıyafet setleri, mont, çizme, kazak, pantolon ve çanta içererek çocuklara umut oluyor. Bağışlarınızla, bu kış da iyilik rüzgarlarını estirerek onların sıcacık gülümsemelerini sağlıyoruz. Bir adet kışlık set bedeli: ₺1.400`,
				Goal:        250000,
				Raised:      175000,
			},
		}
		db.Create(&campaigns)
	}
	var faqCount int64
	db.Model(&model.FAQ{}).Count(&faqCount)
	if faqCount == 0 {
		log.Println("Veritabanı (SSS) boş, başlangıç verileri ekleniyor...")
		faqs := []model.FAQ{
			{Question: "Nasıl bağış yapabilirim?", Answer: "Web sitemizdeki 'Bağış Yap' butonlarını kullanarak online ve güvenli bir şekilde bağış yapabilirsiniz. Ayrıca, banka hesap numaralarımıza EFT/Havale yoluyla da destek olabilirsiniz."},
			{Question: "Bağışlarımın takibini nasıl yapabilirim?", Answer: "Bağışlarınızın ihtiyaç sahiplerine ulaştırılması sürecini sitemizdeki 'Projeler' ve 'Haberler' bölümlerinden şeffaf bir şekilde takip edebilirsiniz."},
			{Question: "Gönüllü olarak nasıl destek olabilirim?", Answer: "'İletişim' sayfamızdaki formu doldurarak veya sosyal medya hesaplarımızdan bize ulaşarak gönüllü ekibimize katılabilirsiniz."},
			{Question: "Yardımlar sadece yurt dışına mı yapılıyor?", Answer: "Hayır, Khayrat AID olarak hem yurt içinde hem de yurt dışında acil insani yardım, eğitim, sağlık ve kalkınma projeleri yürütmekteyiz."},
		}
		db.Create(&faqs)
	
	}
}
	