package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/passageidentity/passage-go"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to read .env variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable required")
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.Handle("/assets/", http.FileServer(http.Dir("./templates")))

	http.ListenAndServe(":"+port, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	inputArgs := map[string]interface{}{"appID": os.Getenv("PASSAGE_APP_ID")}
	outputHTML(w, "templates/index.html", inputArgs)
}

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate this request using the Passage SDK.
	psg, err := passage.New(os.Getenv("PASSAGE_APP_ID"), &passage.Config{APIKey: os.Getenv("PASSAGE_API_KEY")})
	if err != nil {
		fmt.Println("Cannot create psg: ", err)
	}
	_, err = psg.AuthenticateRequest(r)
	if err != nil {
		fmt.Println("Authentication Failed:", err)
		http.ServeFile(w, r, "templates/unauthorized.html")
		return
	}
	inputArgs := map[string]interface{}{"appID": os.Getenv("PASSAGE_APP_ID")}
	outputHTML(w, "templates/dashboard.html", inputArgs)
}
