package ghjournal

import (
	// "net/http"

	el "github.com/deoxxa/echo-logrus"
	"github.com/labstack/echo"
	log "github.com/Sirupsen/logrus"
	mw "github.com/labstack/echo/middleware"
)

func RunServer(port string) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(el.New())
	e.Use(mw.Recover())

	// Routes
	// e.Get("/", hello)

	// Serve static files
	e.Index("www/index.html")
	e.Static("/", "www")

	// Start server
	log.Printf("Starting server on port %s", port)
	e.Run(":" + port)
}

// // Handler
// func hello(c *echo.Context) error {
// 	return c.String(http.StatusOK, "Hello, World!\n")
// }
