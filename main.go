package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	// Start server
	e.Logger.Fatal(e.Start(":1323"))

	// scrapper.Scrape("term")
}

// Handler
func handleHome(c echo.Context) error {
	return c.File("home.html")
}

// Handler
func handleScrape(c echo.Context) error {
	// term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	return nil
	// return c.File("home.html")
}
