package main

import (
	"log"
	"net/http"

	"github.com/9d77v/go-pkg/env"
	"github.com/9d77v/short-url/app"
	"github.com/gin-gonic/gin"
)

type URLConvert struct {
	URL  string `form:"url" json:"url" binding:"required"`
	Hour int    `form:"hour" json:"hour"`
}

var token = env.String("token", "wzCFBZpKqQfVX1T7kL")

func main() {
	err := app.GetDB().AutoMigrate(
		&app.ShortURL{},
	)
	if err != nil {
		log.Println("auto migrate error:", err)
	}

	r := gin.Default()
	r.POST("/url_convert", func(c *gin.Context) {
		if c.Query("token") != token {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "token error",
			})
			return
		}
		var json URLConvert
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		code := app.ConvertURL(json.URL, json.Hour)
		if code == "" {
			c.JSON(200, gin.H{
				"message": "generate failed",
			})
		} else {
			c.JSON(200, gin.H{
				"data": code,
			})
		}
	})
	r.GET("/:short_code", func(c *gin.Context) {
		shortCode := c.Param("short_code")
		url, err := app.GetURL(shortCode)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "404 not found",
			})
		} else {
			c.Redirect(302, url)
		}
	})
	r.Run()
}
