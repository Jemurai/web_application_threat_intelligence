package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type page struct {
	Title    string
	Error    string
	Repsheet bool
}

type recaptchaResponse struct {
	Success            bool
	ChallengeTimestamp string `json:challenge_ts`
	Hostname           string
	ErrorCodes         []string
}

var templates = template.Must(template.ParseFiles("index.html", "admin.html"))

func verifyRecaptcha(gResponse string) bool {
	data := url.Values{}
	data.Set("secret", "6LcuXBAUAAAAAAFcwv--LwXc1mU5C_yYfZICZDCM")
	data.Set("response", gResponse)

	request, _ := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	transport := http.Transport{}
	response, _ := transport.RoundTrip(request)
	defer response.Body.Close()

	bodyString, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var parsedResponse recaptchaResponse
	json.Unmarshal(bodyString, &parsedResponse)

	if !parsedResponse.Success {
		log.Println("Recaptcha validation failed")
	}

	return parsedResponse.Success
}

func repsheetHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["X-Repsheet"] != nil {
			context.Set(r, "repsheet", true)
		}
		next.ServeHTTP(w, r)
	})
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("Content-type", "text/html")
		err := request.ParseForm()
		if err != nil {
			http.Error(response, fmt.Sprintf("error parsing url %v", err), 500)
		}

		p := page{Title: "Login"}

		if context.Get(request, "repsheet") != nil {
			log.Println("MARKED")
			p.Repsheet = true
		}

		templates.ExecuteTemplate(response, "index.html", p)
	} else if request.Method == "POST" {
		request.ParseForm()
		username := request.PostFormValue("inputEmail")
		password := request.PostFormValue("inputPassword")
		recaptcha := request.PostFormValue("g-recaptcha-response")
		recaptchaValid := verifyRecaptcha(recaptcha)

		if username == "admin@example.com" && password == "P4$$w0rd!" && recaptchaValid {
			log.Println("Successfull login for admin@example.com")
			http.Redirect(response, request, "/admin", 302)
		} else {
			log.Println("Login failed for admin@example.com")
			p := page{Title: "Login", Error: "Username or Password Invalid"}
			templates.ExecuteTemplate(response, "index.html", p)
		}
	}
}

func adminHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/html")
	err := request.ParseForm()
	if err != nil {
		http.Error(response, fmt.Sprintf("error parsing url %v", err), 500)
	}
	templates.ExecuteTemplate(response, "admin.html", page{Title: "Admin"})
}

func main() {
	logFile, err := os.OpenFile("logs/app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error accessing log file:", err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.Handle("/", handlers.ProxyHeaders(handlers.LoggingHandler(logFile, repsheetHandler(http.HandlerFunc(loginHandler)))))
	r.Handle("/admin", handlers.ProxyHeaders(handlers.LoggingHandler(logFile, repsheetHandler(http.HandlerFunc(adminHandler)))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", r)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
