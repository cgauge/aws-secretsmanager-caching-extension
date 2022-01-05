package server

import (
	"net/http"

	"github.com/cgauge/aws-secretsmanager-caching-extension/secrets"
)

func StartHTTPServer(port string) {
	println(secrets.PrintPrefix, "Starting Httpserver on port ", port)

	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		refresh := r.URL.Query().Get("refresh")

		value := secrets.GetSecretCache(name, refresh)

		if len(value) != 0 {
			_, _ = w.Write([]byte(value))
		} else {
			_, _ = w.Write([]byte("No data found"))
		}
	})

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		panic(err)
	}
}
