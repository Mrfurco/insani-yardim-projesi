package handler

import (
	"net/http"
	"github.com/MrFurco/insani-yardim-projesi/internal/model"
	"github.com/MrFurco/insani-yardim-projesi/internal/repository"
	"github.com/gin-gonic/gin" // Gin importunu da ekleyerek tutarlı hale getiriyoruz.
)

// ListPosts, tüm haber yazılarını veritabanından çeker ve listeler.
func ListPosts(c *gin.Context) {
	var posts []model.Post
	repository.DB.Order("id desc").Find(&posts)

	c.HTML(http.StatusOK, "admin_list_posts.html", gin.H{
		"title": "Haber Yönetimi",
		"data":  gin.H{"posts": posts},
	})
}

// ShowNewPostForm, yeni haber oluşturma formunu gösterir.
func ShowNewPostForm(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_new_post.html", gin.H{
		"title": "Yeni Haber Ekle",
	})
}

// CreatePost, formdan gelen verilerle yeni bir haber oluşturur.
func CreatePost(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")

	post := model.Post{
		Title:   title,
		Content: content,
	}

	if result := repository.DB.Create(&post); result.Error != nil {
		c.String(http.StatusInternalServerError, "Hata oluştu: %v", result.Error)
		return
	}

	c.Redirect(http.StatusFound, "/admin/haberler")
}