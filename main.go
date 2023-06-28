package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type PhoneData struct {
	Phone string
	Name  string
}

var phoneData []PhoneData

func main() {
	fs := http.FileServer(http.Dir("./usage"))
	http.Handle("/usage/", http.StripPrefix("/usage/", fs))
	data, err := ioutil.ReadFile("info.txt")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			phone := strings.TrimSpace(parts[0])
			name := strings.TrimSpace(parts[1])
			phoneData = append(phoneData, PhoneData{Phone: phone, Name: name})
		}
	}

	
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.ListenAndServe(":8000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("form.html"))
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	phoneNumber := r.Form.Get("phone")
	if string(phoneNumber[0]) == "8" {
		phoneNumber = "+7" + phoneNumber[1:]
	}
	fmt.Println(phoneNumber)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	for _, data := range phoneData {
		if data.Phone == phoneNumber {
			tmpl := template.Must(template.ParseFiles("home.html"))
			tmpl.Execute(w, data.Name)
			return
		}
	}

	tmpl := template.Must(template.ParseFiles("error.html"))
	tmpl.Execute(w, nil)
}
