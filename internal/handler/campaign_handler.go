package handler

import (
	"net/http"
	"strconv"

	"github.com/MrFurco/insani-yardim-projesi/internal/model"
	"github.com/MrFurco/insani-yardim-projesi/internal/repository"
	"github.com/gin-gonic/gin"
)

// --- PUBLIC HANDLERS ---

// ListCampaigns, ana sayfa için gerekli tüm verileri çeker ve gönderir.
func ListCampaigns(c *gin.Context) {
	// Kampanyaları çekiyoruz
	var campaigns []model.Campaign
	repository.DB.Order("id asc").Find(&campaigns)

	// Kampanyalar için View Model oluşturuyoruz (Yüzde hesabı için)
	type CampaignView struct {
		ID, Goal, Raised, Percentage int
		Title                        string
	}
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
			ID: campaign.ID, Title: campaign.Title, Goal: campaign.Goal, Raised: campaign.Raised, Percentage: percentage,
		})
	}

	// Tüm verileri tek bir map içinde HTML'e gönderiyoruz
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Ana Sayfa",
		"data": gin.H{
			"campaigns": campaignViews,
			// Diğer bölümleri eklediğimizde buraya "posts", "faqs" gibi veriler de gelecek
		},
	})
}

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
		"data":  gin.H{"ID": campaign.ID, "Title": campaign.Title, "Description": campaign.Description, "Goal": campaign.Goal, "Raised": campaign.Raised, "Percentage": percentage},
	})
}

// --- ADMIN HANDLERS ---

func ShowAdminDashboard(c *gin.Context) {
	var campaigns []model.Campaign
	repository.DB.Order("id asc").Find(&campaigns)
	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"title": "Admin Paneli",
		"data":  gin.H{"campaigns": campaigns},
	})
}

func ShowNewCampaignForm(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_new_campaign.html", gin.H{
		"title": "Yeni Kampanya Ekle",
	})
}

func CreateCampaign(c *gin.Context) {
	title := c.PostForm("title")
	goal, _ := strconv.Atoi(c.PostForm("goal"))
	// Not: Henüz formda description alanı yok, bir sonraki adımda ekleyeceğiz.
	campaign := model.Campaign{Title: title, Goal: goal, Raised: 0}
	repository.DB.Create(&campaign)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func ShowEditCampaignForm(c *gin.Context) {
	id := c.Param("id")
	var campaign model.Campaign
	if err := repository.DB.First(&campaign, id).Error; err != nil {
		c.String(http.StatusNotFound, "Düzenlenecek kampanya bulunamadı.")
		return
	}
	c.HTML(http.StatusOK, "admin_edit_campaign.html", gin.H{
		"title": "Kampanya Düzenle",
		"data":  campaign,
	})
}

func UpdateCampaign(c *gin.Context) {
	id := c.Param("id")
	var campaign model.Campaign
	if err := repository.DB.First(&campaign, id).Error; err != nil {
		c.String(http.StatusNotFound, "Güncellenecek kampanya bulunamadı.")
		return
	}
	campaign.Title = c.PostForm("title")
	campaign.Goal, _ = strconv.Atoi(c.PostForm("goal"))
	campaign.Raised, _ = strconv.Atoi(c.PostForm("raised"))
	// Not: Henüz formda description alanı yok, bir sonraki adımda ekleyeceğiz.
	repository.DB.Save(&campaign)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func DeleteCampaign(c *gin.Context) {
	id := c.Param("id")
	repository.DB.Delete(&model.Campaign{}, id)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}