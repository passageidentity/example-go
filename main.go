package main

import (
	"net/http"

	"github.com/passageidentity/passage-go"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	psg := passage.New("demo")
	_, err := psg.AuthenticateRequest(r)
	if err != nil {
		http.ServeFile(w, r, "templates/unauthorized.html")
		return
	}

	http.ServeFile(w, r, "templates/dashboard.html")
}
