package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vitezprchal/waitlistgo/internal/db"
	"github.com/vitezprchal/waitlistgo/internal/forms"
	"github.com/vitezprchal/waitlistgo/internal/models"
	"github.com/vitezprchal/waitlistgo/internal/renderer"
	"github.com/vitezprchal/waitlistgo/website/view"
)

func InitServer(port string) {

	repository := db.InitDatabase()

	if port == "" {
		fmt.Println("Port not provided, using default port 8080")
		port = ":8080"
	}

	router := gin.Default()

	router.LoadHTMLFiles("./view/*.html")

	router.SetTrustedProxies(nil)

	gin_html := router.HTMLRender
	router.HTMLRender = &renderer.HTMLTemplRenderer{FallbackHtmlRenderer: gin_html}

	var home_seo = &models.SEO{
		Title:       "WaitlistGo",
		Description: "WaitlistGo is a simple waitlist management system.",
		Keywords:    "waitlist, management, system",
		AuthorName:  "Vitezslav",
		ImageURL:    "https://www.something.some/assets/logo-image.png",
		PageURL:     "https://www.something.some",
	}

	router.StaticFile("./build/styles.css", "./build/styles.css")
	router.StaticFile("./build/main.js", "./build/main.js")

	router.Static("/assets", "./website/assets")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", view.Home(home_seo))
	})

	router.POST("/submit", func(c *gin.Context) {
		var form forms.EmailForm

		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if form.Terms != "on" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You must agree to the terms"})
			return
		}

		err := repository.AddSubscriber(c, &models.Subscriber{
			Name:  form.Name,
			Email: form.Email,
		})

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error adding email to the waitlist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email added to the waitlist"})
	})

	router.Run(port)

	defer repository.Close()
}
