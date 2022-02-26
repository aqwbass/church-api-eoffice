package main

import (
	"fmt"
	"net/http"
	"os"

	myMiddL "github.com/church/church-api-eoffice/middleware"
	"github.com/church/church-api-eoffice/route"
	"github.com/labstack/echo/v4"
	echoMiddL "github.com/labstack/echo/v4/middleware"
)

func GetENV(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

var (
	APP_PORT          = GetENV("APP_PORT", "3000")
	PSQL_DATABASE_URL = GetENV("PSQL_DATABASE_URL", "postgres://postgres:postgres@psql_db:5432/app_example?sslmode=disable")
)

func main() {
	e := echo.New()
	e.Use(echoMiddL.Logger())
	e.Use(echoMiddL.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	middL := myMiddL.InitMiddleware("test")
	e.Use(echoMiddL.Recover())
	e.Use(echoMiddL.CORSWithConfig(echoMiddL.CORSConfig{
		Skipper:      echoMiddL.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middL.InitContextIfNotExists)
	e.Use(middL.InputForm)

	r := route.NewRoute(e, middL)

	/* repository */

	/* usecase */

	/* handler */

	/* validate */

	/* inject route */

	fmt.Println(r)
	/* serve echo */
	port := fmt.Sprintf(":%s", APP_PORT)
	e.Logger.Fatal(e.Start(port))

}
