package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Shortened struct {
	gorm.Model
	Name string `json:"name"`
	Url  string `json:"url"`
}

func randomName() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func main() {
	router := mux.NewRouter()

	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(Shortened{})

	router.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		short := Shortened{}
		json.NewDecoder(r.Body).Decode(&short)
		if short.Url == "" {
			w.WriteHeader(500)
			return
		}
		if _, err := url.Parse(short.Url); err != nil {
			w.WriteHeader(500)
			return
		}

		if short.Name != "" {
			existingShort := Shortened{}
			err := db.Where("name = ?", short.Name).First(&existingShort).Error
			if err != nil {
				db.Create(&short)
				json.NewEncoder(w).Encode(short)
			} else {
				json.NewEncoder(w).Encode(existingShort)
			}
			return
		}
		for {
			rName := randomName()
			existingShort := Shortened{}
			err := db.Where("name = ?", rName).First(&existingShort).Error
			if err != nil {
				short.Name = rName
				db.Create(&short)
				json.NewEncoder(w).Encode(short)
				break
			}
		}

	}).Methods("POST")

	router.HandleFunc("/r/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name, ok := vars["name"]
		if !ok {
			w.WriteHeader(500)
			return
		}
		short := Shortened{}
		if err := db.Where("name = ?", name).First(&short).Error; err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Location", short.Url)
		w.WriteHeader(302)
	}).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("webbone/dist")))

	http.ListenAndServe(":8080", router)
}
