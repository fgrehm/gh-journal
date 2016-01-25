package ghjournal

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	el "github.com/deoxxa/echo-logrus"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

func RunServer(port string) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(el.New())
	e.Use(mw.Recover())

	// Routes
	e.Get("/report/:date", showReport)

	// Serve static files
	e.Index("www/index.html")
	e.Static("/", "www")

	// Start server
	log.Printf("Starting server on port %s", port)
	e.Run(":" + port)
}

func showReport(c *echo.Context) error {
	date, err := time.Parse(time.RFC3339, c.Param("date")+"T00:00:00Z")
	if err != nil {
		return err
	}

	data, err := buildReport(date)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}
