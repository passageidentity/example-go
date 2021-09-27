package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/passageidentity/passage-go"
)

func main() {
	os.Setenv("PASSAGE_APP_ID", "[YOUR_APP_ID_HERE]")
	os.Setenv("PASSAGE_API_KEY", "[YOUR_PASSAGE_API_KEY_HERE]")
	os.Setenv("PORT", "5000")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.Handle("/assets/", http.FileServer(http.Dir("./templates")))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
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
	userID, err := psg.AuthenticateRequest(r)
	if err != nil {
		fmt.Println("Authentication Failed:", err)
		http.ServeFile(w, r, "templates/unauthorized.html")
		return
	}
	user, err := psg.GetUser(userID)
	if err != nil {
		fmt.Println("Could not get user: ", err)
		return
	}
	inputArgs := map[string]interface{}{"email": user.Email}
	outputHTML(w, "templates/dashboard.html", inputArgs)
}
