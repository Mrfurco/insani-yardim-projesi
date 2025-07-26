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
	router.Static("/static", "./static")

	// Şablonları en basit ve güvenilir şekilde yüklüyoruz.
	router.LoadHTMLGlob("templates/*.html")

	// --- Public Rotalar ---
	router.GET("/", handler.ListCampaigns)
	router.GET("/kampanya/:id", handler.ShowCampaign)
	router.GET("/kampanyalar", handler.ShowCampaignsPage)
	router.GET("/haberler", handler.ShowNewsPage)
	router.GET("/sss", handler.ShowFaqsPage)
	router.GET("/medya", handler.ShowMediaPage)
	router.GET("/hakkimizda", handler.ShowAboutPage)
	router.GET("/iletisim", handler.ShowContactPage)

	// --- Admin Grubu ve Rotaları ---
	adminGroup := router.Group("/admin")
	adminGroup.Use(gin.BasicAuth(gin.Accounts{
		"admin": "12345",
	}))
	// Kampanya Yönetimi
	adminGroup.GET("/dashboard", handler.ShowAdminDashboard)
	adminGroup.GET("/kampanyalar/yeni", handler.ShowNewCampaignForm)
	adminGroup.POST("/kampanyalar", handler.CreateCampaign)
	adminGroup.GET("/kampanyalar", handler.ListCampaignsAdmin)
	adminGroup.GET("/kampanyalar/sil/:id", handler.DeleteCampaign)
	adminGroup.GET("/kampanyalar/duzenle/:id", handler.ShowEditCampaignForm)
	adminGroup.POST("/kampanyalar/duzenle/:id", handler.UpdateCampaign)
	// Haber Yönetimi
	adminGroup.GET("/haberler", handler.ListPosts)
	adminGroup.GET("/haberler/yeni", handler.ShowNewPostForm)
	adminGroup.POST("/haberler", handler.CreatePost)
	adminGroup.GET("/haberler/duzenle/:id", handler.ShowEditPostForm)
	adminGroup.POST("/haberler/duzenle/:id", handler.UpdatePost)
	adminGroup.GET("/haberler/sil/:id", handler.DeletePost)

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
				Title:       "Bu Kış Kimse Üşümesin",
				Description: `Yurt içinde ve yurt dışında, kışın ağır koşullarında yaşam mücadelesi veren birçok aile, soğuktan korunmakta zorluk çekiyor...`,
				ImageURL:    "https://images.pexels.com/photos/14028471/pexels-photo-14028471.jpeg?auto=compress&cs=tinysrgb&w=600",
				Goal:        250000,
				Raised:      175000,
			},
			{
				Title:       "Filistin İçin Acil Yardım",
				Description: `Gazze'de yaşanan insanlık dramına sessiz kalmayın. Acil gıda, su ve tıbbi malzeme yardımlarınızla bir hayata umut olabilirsiniz.`,
				ImageURL:    "https://images.pexels.com/photos/18451314/pexels-photo-18451314.jpeg?auto=compress&cs=tinysrgb&w=600",
				Goal:        500000,
				Raised:      95000,
			},
			{
				Title:       "Bir Burs Bir Gelecek",
				Description: `Geleceğimizin güvence altına almanın en güzel yolu eğitimdir. Öğrencilerimizin eğitim hayatlarına devam edebilmeleri için burs desteği sağlıyoruz.`,
				ImageURL:    "https://images.pexels.com/photos/267885/pexels-photo-267885.jpeg?auto=compress&cs=tinysrgb&w=600",
				Goal:        100000,
				Raised:      85000,
			},
			{
				Title:       "Temiz Su Hayattır",
				Description: `Afrika'da suya erişimi olmayan köylerde su kuyuları açarak binlerce insana temiz su ve sağlıklı bir yaşam sunuyoruz.`,
				ImageURL:    "https://images.pexels.com/photos/6646905/pexels-photo-6646905.jpeg?auto=compress&cs=tinysrgb&w=600",
				Goal:        100000,
				Raised:      60000,
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
