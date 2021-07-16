package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/passageidentity/passage-go"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.Handle("/assets/", http.FileServer(http.Dir("./templates")))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate this request using the Passage SDK.
	psg, err := passage.New("UKbRUx")
	if err != nil {
		fmt.Println("Cannot create psg: ", err)
	}
	_, errAuth := psg.AuthenticateRequest(r)
	if errAuth != nil {
		fmt.Println("Authentication Failed:", err)
		http.ServeFile(w, r, "templates/unauthorized.html")
		return
	}

	http.ServeFile(w, r, "templates/dashboard.html")
}
