package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
)

func main() {
	r := gin.Default()

	tmpl, err := template.New("").ParseFiles("views/index.html")
	if err != nil {
		log.Fatalf("Şablon yüklenirken hata oluştu: %v", err)
	}
	r.SetHTMLTemplate(tmpl)

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Ana Sayfa",  // Dinamik veri
		})
	})
	r.GET("/about", func(c *gin.Context) {
		c.HTML(200, "about.html", gin.H{
			"title": "Ana Sayfa",  // Dinamik veri
		})
	})
	r.GET("/contact", func(c *gin.Context) {
		c.HTML(200, "contact.html", gin.H{
			"title": "Ana Sayfa",  // Dinamik veri
		})
	})

	// Sunucuyu başlat
	log.Fatal(r.Run(":8080"))
}

