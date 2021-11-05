package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/golang-jwt/jwt"
)

var publicKey *rsa.PublicKey

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable required")
	}

	// Get and parse public key
	// NOTE: THIS IS REQUIRED FOR RUNNING IN NON-PROD, BECAUSE THE SDK IS ONLY SET
	// UP TO WORK IN PRODUCTION
	psg_public_key := os.Getenv("PASSAGE_PUBLIC_KEY")
	if port == "" {
		log.Fatal("PASSAGE_PUBLIC_KEY environment variable required")
	}

	// Parse the returned public key string to an rsa.PublicKey:
	publicKeyBytes, err := base64.RawURLEncoding.DecodeString(psg_public_key)
	if err != nil {
		log.Fatal("could not parse Passage App's public key: expected valid base-64")
	}
	pemBlock, _ := pem.Decode(publicKeyBytes)
	if pemBlock == nil {
		log.Fatal("could not parse Passage App's public key: missing PEM block")
	}
	publicKey, err = x509.ParsePKCS1PublicKey(pemBlock.Bytes)
	if err != nil {
		log.Fatal("could not parse Passage App's public key: invalid PKCS #1 public key")
	}
	// END WEIRD PUBLIC KEY STUFF

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
	// Manually validate JWTs with Public Key
	userID, err := authenticateRequestWithCookie(r)
	if err != nil {
		fmt.Println("Authentication Failed:", err)
		http.ServeFile(w, r, "templates/unauthorized.html")
		return
	}

	// PRINT USER ID INSTEAD OF EMAIL IN UAT BECAUSE WE CAN'T USE THE SDK AND IM LAZY
	inputArgs := map[string]interface{}{"email": userID}
	outputHTML(w, "templates/dashboard.html", inputArgs)
}

// AuthenticateRequestWithCookie fetches a cookie from the request and uses it to authenticate
// returns the userID (string) on success, error on failure
func authenticateRequestWithCookie(r *http.Request) (string, error) {
	authTokenCookie, err := r.Cookie("psg_auth_token")
	if err != nil {
		return "", errors.New("missing authentication token: expected \"psg_auth_token\" cookie")
	}

	userID, valid := validateAuthToken(authTokenCookie.Value)
	if !valid {
		return "", errors.New("invalid authentication token")
	}

	return userID, nil
}

// ValidateAuthToken determines whether a JWT is valid or not
// returns userID (string) on success, error on failure
func validateAuthToken(authToken string) (string, bool) {
	// Verify that the authentication token is valid:
	parsedToken, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing algorithm")
		}
		return publicKey, nil
	})
	if err != nil {
		return "", false
	}

	// Extract claims from JWT:
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", false
	}

	return userID, true
}
