package handler

import (
	"net/http"
	"strconv"

	"github.com/MrFurco/insani-yardim-projesi/internal/model"
	"github.com/MrFurco/insani-yardim-projesi/internal/repository"
	"github.com/MrFurco/insani-yardim-projesi/internal/session"
	"github.com/gin-gonic/gin"
)

// AddToCart, sepete yeni bir ürün ekler.
func AddToCart(c *gin.Context) {
	session, _ := session.Store.Get(c.Request, "khayrat-session")

	campaignIDStr := c.PostForm("campaignId")
	amountStr := c.PostForm("amount")
	itemName := c.PostForm("itemName")

	campaignID, _ := strconv.ParseUint(campaignIDStr, 10, 32)
	amount, _ := strconv.Atoi(amountStr)

	newItem := model.CartItem{
		CampaignID: uint(campaignID),
		ItemName:   itemName,
		Amount:     amount,
	}

	var cart []model.CartItem
	if session.Values["cart"] != nil {
		// Tip dönüşümü yaparken hata kontrolü daha güvenlidir, ama şimdilik böyle devam edelim
		cart = session.Values["cart"].([]model.CartItem)
	}

	cart = append(cart, newItem)
	session.Values["cart"] = cart
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/sepet")
}

// ShowCartPage, kullanıcının sepetindeki ürünleri gösterir.
func ShowCartPage(c *gin.Context) {
	session, _ := session.Store.Get(c.Request, "khayrat-session")
	var cart []model.CartItem
	if val, ok := session.Values["cart"]; ok && val != nil {
		cart = val.([]model.CartItem)
	}

	c.HTML(http.StatusOK, "sepet.html", gin.H{
		"title": "Bağış Sepeti",
		"cart":  cart,
	})
}

// ShowCheckoutPage, ödeme formunu ve sepet özetini gösterir.
func ShowCheckoutPage(c *gin.Context) {
	session, _ := session.Store.Get(c.Request, "khayrat-session")
	var cart []model.CartItem
	if val, ok := session.Values["cart"]; ok && val != nil {
		cart = val.([]model.CartItem)
	}

	if len(cart) == 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var total int
	for _, item := range cart {
		total += item.Amount
	}

	c.HTML(http.StatusOK, "odeme.html", gin.H{
		"title": "Bağışı Tamamla",
		"total": total,
	})
}

// ProcessDonation, ödeme formundan gelen bilgileri ve sepetteki verileri işler.
func ProcessDonation(c *gin.Context) {
	session, _ := session.Store.Get(c.Request, "khayrat-session")
	var cart []model.CartItem
	if val, ok := session.Values["cart"]; ok && val != nil {
		cart = val.([]model.CartItem)
	}

	if len(cart) == 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Her bir sepet kalemi için ayrı bir bağış kaydı oluştur ve kampanyayı güncelle
	for _, item := range cart {
		donation := model.Donation{
			DonorName:  c.PostForm("name"),
			DonorEmail: c.PostForm("email"),
			Amount:     item.Amount,
			CampaignID: item.CampaignID,
		}
		repository.DB.Create(&donation)

		var campaign model.Campaign
		repository.DB.First(&campaign, item.CampaignID)
		campaign.Raised += item.Amount
		repository.DB.Save(&campaign)
	}

	session.Values["cart"] = nil
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/tesekkurler")
}

// ShowThankYouPage, basit bir teşekkür sayfası gösterir.
func ShowThankYouPage(c *gin.Context) {
	c.HTML(http.StatusOK, "tesekkurler.html", gin.H{
		"title": "Teşekkür Ederiz",
	})
}