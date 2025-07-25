package handler

import (
	"html/template"
	"net/http"

	"github.com/MrFurco/insani-yardim-projesi/internal/model"
	"github.com/MrFurco/insani-yardim-projesi/internal/repository"
	"github.com/gin-gonic/gin"
)

// ListPosts, tüm haber yazılarını veritabanından çeker ve listeler.
func ListPosts(c *gin.Context) {
	var posts []model.Post
	repository.DB.Order("id desc").Find(&posts)

	tmpl, err := template.ParseFiles("templates/admin_layout.html", "templates/admin_list_posts.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parse error: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"title": "Haber Yönetimi",
		"data":  gin.H{"posts": posts},
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execute error: %v", err)
		return
	}
}

// ShowNewPostForm, yeni haber oluşturma formunu gösterir.
func ShowNewPostForm(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/admin_layout.html", "templates/admin_new_post.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parse error: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"title": "Yeni Haber Ekle",
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execute error: %v", err)
		return
	}
}

// CreatePost, formdan gelen verilerle yeni bir haber oluşturur.
func CreatePost(c *gin.Context) {
	// Yeni alanları formdan alıyoruz
	title := c.PostForm("title")
	summary := c.PostForm("summary")
	content := c.PostForm("content")
	imageURL := c.PostForm("image_url")

	post := model.Post{
		Title:    title,
		Summary:  summary,
		Content:  content,
		ImageURL: imageURL,
	}

	if result := repository.DB.Create(&post); result.Error != nil {
		c.String(http.StatusInternalServerError, "Hata oluştu: %v", result.Error)
		return
	}

	c.Redirect(http.StatusFound, "/admin/haberler")
}

// --- YENİ EKLENEN FONKSİYONLAR ---

// ShowEditPostForm, belirli bir haberi düzenleme formunu gösterir.
func ShowEditPostForm(c *gin.Context) {
	id := c.Param("id")
	var post model.Post
	if err := repository.DB.First(&post, id).Error; err != nil {
		c.String(http.StatusNotFound, "Düzenlenecek haber bulunamadı.")
		return
	}

	tmpl, err := template.ParseFiles("templates/admin_layout.html", "templates/admin_edit_post.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parse error: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
		"title": "Haberi Düzenle",
		"data":  post,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execute error: %v", err)
		return
	}
}

// UpdatePost, formdan gelen verilerle belirli bir haberi günceller.
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post model.Post
	if err := repository.DB.First(&post, id).Error; err != nil {
		c.String(http.StatusNotFound, "Güncellenecek haber bulunamadı.")
		return
	}

	// Yeni alanları formdan alıp güncelliyoruz
	post.Title = c.PostForm("title")
	post.Summary = c.PostForm("summary")
	post.Content = c.PostForm("content")
	post.ImageURL = c.PostForm("image_url")

	repository.DB.Save(&post)
	c.Redirect(http.StatusFound, "/admin/haberler")
}

// DeletePost, belirli bir haberi siler.
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	repository.DB.Delete(&model.Post{}, id)
	c.Redirect(http.StatusFound, "/admin/haberler")
}
