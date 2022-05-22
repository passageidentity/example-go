package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/passageidentity/passage-go"
)

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		psg, _ := passage.New("<Passage App ID>", nil)
		_, err := psg.AuthenticateRequest(c.Request)
		if err != nil {
			// Authentication failed!
			c.HTML(http.StatusForbidden, "unauthorized.html", nil)
			c.Abort()
		}

		// Authentication was successful, proceed.
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("html/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	authenticated := r.Group("/", authRequired())

	authenticated.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", nil)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
