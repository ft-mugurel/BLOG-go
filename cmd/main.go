package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
)

func main() {
	r := gin.Default()

	// Tüm HTML dosyalarını yüklüyoruz
	tmpl, err := template.ParseFiles(
		"views/index.html",
		"views/about.html",
		"views/contact.html",
	)
	if err != nil {
		log.Fatalf("Şablonlar yüklenirken hata oluştu: %v", err)
	}
	r.SetHTMLTemplate(tmpl)

	// Sayfalar
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Ana Sayfa",
		})
	})
	r.GET("/about", func(c *gin.Context) {
		c.HTML(200, "about.html", gin.H{
			"title": "Hakkımda",
		})
	})
	r.GET("/contact", func(c *gin.Context) {
		c.HTML(200, "contact.html", gin.H{
			"title": "İletişim",
		})
	})

	log.Fatal(r.Run(":8080"))
}

