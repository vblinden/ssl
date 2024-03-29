package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/asdgo/asdgo"
)

func main() {
	asdgo.New(asdgo.Config{
		CsrfExempts: []string{"/"},
	})

	asdgo.Router().Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, strings.Trim(`
microservice to get x509 certificate info

post a x509 certificate in pem format to the endpoint

example:
  $ curl -X POST -d "$(cat cert.pem)" https://ssl.vblinden.dev
  { "subject": "...", "issuer": "...", "not_before": "...", "not_after": "..."}

		`, "\n"))
	})

	asdgo.Router().Post("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading request body", http.StatusInternalServerError)
			return
		}
		block, _ := pem.Decode(body)
		if block == nil {
			http.Error(w, "error parsing certificate", http.StatusBadRequest)
			return
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing certificate: %s", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(
			w,
			`{"subject": "%s","issuer": "%s","not_before": "%s","not_after": "%s"}`,
			cert.Subject.CommonName,
			cert.Issuer.CommonName,
			cert.NotBefore,
			cert.NotAfter,
		)
	})

	fmt.Println("starting server on :3000")
	http.ListenAndServe(":3000", asdgo.Router())
}
