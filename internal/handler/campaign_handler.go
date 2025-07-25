package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/MrFurco/insani-yardim-projesi/internal/model"
	"github.com/MrFurco/insani-yardim-projesi/internal/repository"
	"github.com/gin-gonic/gin"
)

// ListCampaigns, ana sayfa için gerekli tüm verileri çeker ve gönderir.
func ListCampaigns(c *gin.Context) {
	// Kampanyaları çekiyoruz
	var campaigns []model.Campaign
	repository.DB.Order("id asc").Find(&campaigns)

	// Kampanyalar için View Model'i fonksiyonun İÇİNDE tanımlıyoruz
	type CampaignView struct {
		ID, Goal, Raised, Percentage int
		Title, Description, ImageURL string
	}

	// Değişkeni de fonksiyonun İÇİNDE tanımlıyoruz
	var campaignViews []CampaignView

	for _, campaign := range campaigns {
		percentage := 0
		if campaign.Goal > 0 {
			percentage = (campaign.Raised * 100) / campaign.Goal
		}
		if percentage > 100 {
			percentage = 100
		}
		campaignViews = append(campaignViews, CampaignView{
			ID:          campaign.ID,
			Title:       campaign.Title,
			Description: campaign.Description,
			ImageURL:    campaign.ImageURL,
			Goal:        campaign.Goal,
			Raised:      campaign.Raised,
			Percentage:  percentage,
		})
	}

	// Haberleri ve SSS'leri de çekiyoruz
	var posts []model.Post
	repository.DB.Order("id desc").Limit(5).Find(&posts)
	var faqs []model.FAQ
	repository.DB.Order("id asc").Find(&faqs)

	// Tüm verileri tek bir map içinde HTML'e gönderiyoruz
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Ana Sayfa",
		"data": gin.H{
			"campaigns": campaignViews,
			"posts":     posts,
			"faqs":      faqs,
		},
	})

} // <-- ListCampaigns fonksiyonu burada bitiyor, fazladan parantez yok.

// ... (Diğer handler fonksiyonları: ShowCampaign, CreateCampaign vb. aynı kalabilir) ...
func ShowCampaign(c *gin.Context) {
	id := c.Param("id")
	var campaign model.Campaign
	if err := repository.DB.First(&campaign, id).Error; err != nil {
		c.String(http.StatusNotFound, "Aradığınız kampanya bulunamadı.")
		return
	}

	percentage := 0
	if campaign.Goal > 0 {
		percentage = (campaign.Raised * 100) / campaign.Goal
	}
	if percentage > 100 {
		percentage = 100
	}

	c.HTML(http.StatusOK, "campaign_detail.html", gin.H{
		"title": campaign.Title,
		"data":  gin.H{"ID": campaign.ID, "Title": campaign.Title, "Description": campaign.Description, "ImageURL": campaign.ImageURL, "Goal": campaign.Goal, "Raised": campaign.Raised, "Percentage": percentage},
	})
}

func ShowAdminDashboard(c *gin.Context) {
	var campaigns []model.Campaign
	repository.DB.Order("id asc").Find(&campaigns)

	var totalGoal int
	var totalRaised int
	for _, campaign := range campaigns {
		totalGoal += campaign.Goal
		totalRaised += campaign.Raised
	}
	totalCampaigns := len(campaigns)

	// Şablonları manuel olarak birleştiriyoruz: Önce layout, sonra içerik.
	tmpl, err := template.ParseFiles("templates/admin_layout.html", "templates/admin_dashboard.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parse error: %v", err)
		return
	}

	// Birleştirdiğimiz şablonlardan "layout" olanı çalıştırıyoruz.
	err = tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"title": "Yönetim Paneli",
		"data": gin.H{
			"TotalCampaigns": totalCampaigns,
			"TotalGoal":      totalGoal,
			"TotalRaised":    totalRaised,
		},
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execute error: %v", err)
		return
	}
}

func ShowNewCampaignForm(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/admin_layout.html", "templates/admin_new_campaign.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parse error: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"title": "Yeni Kampanya Ekle",
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execute error: %v", err)
		return
	}
}

func CreateCampaign(c *gin.Context) {
	var imageURLForDB string // Veritabanına kaydedilecek nihai URL

	// 1. Dosya yüklenmiş mi diye kontrol et
	file, err := c.FormFile("image_file")
	if err == nil {
		// Dosya yüklenmişse, onu sunucuya kaydet
		// Benzersiz bir dosya adı oluştur (zaman damgası + orijinal ad)

		newFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filePath := filepath.Join("static", "uploads", newFileName)

		// Dosyayı kaydet
		if err := c.SaveUploadedFile(file, filePath); err == nil {
			// Başarılı olursa, veritabanına kaydedilecek URL'i ayarla
			imageURLForDB = "/" + filepath.ToSlash(filePath)
		}
	}

	// 2. Dosya yüklenmemişse veya hata oluşmuşsa, URL alanına bak
	if imageURLForDB == "" {
		imageURLForDB = c.PostForm("image_url")
	}

	// 3. Veritabanına kaydet
	title := c.PostForm("title")
	description := c.PostForm("description")
	goal, _ := strconv.Atoi(c.PostForm("goal"))

	campaign := model.Campaign{Title: title, Description: description, ImageURL: imageURLForDB, Goal: goal, Raised: 0}
	repository.DB.Create(&campaign)
	c.Redirect(http.StatusFound, "/admin/kampanyalar")
}

func ShowEditCampaignForm(c *gin.Context) {
	id := c.Param("id")
	var campaign model.Campaign
	if err := repository.DB.First(&campaign, id).Error; err != nil {
		c.String(http.StatusNotFound, "Düzenlenecek kampanya bulunamadı.")
		return
	}

	tmpl, err := template.ParseFiles("templates/admin_layout.html", "templates/admin_edit_campaign.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parse error: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"title": "Kampanya Düzenle",
		"data":  campaign,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execute error: %v", err)
		return
	}
}

func UpdateCampaign(c *gin.Context) {
	id := c.Param("id")
	var campaign model.Campaign
	if err := repository.DB.First(&campaign, id).Error; err != nil {
		c.String(http.StatusNotFound, "Güncellenecek kampanya bulunamadı.")
		return
	}

	var imageURLForDB string
	file, err := c.FormFile("image_file")
	if err == nil {
		// Yeni bir dosya yüklendiyse, eskisini silip yenisini kaydet (şimdilik silme adımını atlayalım)

		newFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filePath := filepath.Join("static", "uploads", newFileName)
		if err := c.SaveUploadedFile(file, filePath); err == nil {
			imageURLForDB = "/" + filepath.ToSlash(filePath)
		}
	}

	if imageURLForDB == "" {
		imageURLForDB = c.PostForm("image_url")
	}

	// Veritabanı objesini güncelle
	campaign.Title = c.PostForm("title")
	campaign.Description = c.PostForm("description")
	campaign.ImageURL = imageURLForDB
	campaign.Goal, _ = strconv.Atoi(c.PostForm("goal"))
	campaign.Raised, _ = strconv.Atoi(c.PostForm("raised"))

	repository.DB.Save(&campaign)
	c.Redirect(http.StatusFound, "/admin/kampanyalar")
}

func DeleteCampaign(c *gin.Context) {
	id := c.Param("id")
	repository.DB.Delete(&model.Campaign{}, id)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// ListCampaignsAdmin, sadece kampanya listesini yönetmek için kullanılır.
func ListCampaignsAdmin(c *gin.Context) {
	var campaigns []model.Campaign
	repository.DB.Order("id asc").Find(&campaigns)

	tmpl, err := template.ParseFiles("templates/admin_layout.html", "templates/admin_list_campaigns.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parse error: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"title": "Kampanya Yönetimi",
		"data":  gin.H{"campaigns": campaigns},
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execute error: %v", err)
		return
	}
}
