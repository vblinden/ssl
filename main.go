package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/asdgo/asdgo"
	"github.com/labstack/echo/v4"
)

func main() {
	asd := asdgo.New(&asdgo.Config{})

	asd.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, strings.Trim(`
microservice to get x509 certificate info

post a x509 certificate in pem format to the endpoint

example:
  $ curl -X POST -d "$(cat cert.pem)" https://ssl.vblinden.dev
  { "subject": "...", "issuer": "...", "not_before": "...", "not_after": "..."}

		`, "\n"))
	})

	asd.POST("/", func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)

		if err != nil {
			return c.String(http.StatusInternalServerError, "error reading request body")
		}

		block, _ := pem.Decode(body)
		if block == nil {
			return c.String(http.StatusBadRequest, "error parsing certificate")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("error parsing certificate: %s", err))
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		return c.String(http.StatusOK, fmt.Sprintf(
			`{"subject": "%s","issuer": "%s","not_before": "%s","not_after": "%s"}`,
			cert.Subject.CommonName,
			cert.Issuer.CommonName,
			cert.NotBefore,
			cert.NotAfter,
		))
	})

	asd.Logger.Fatal(asd.Start(":3000"))
}
