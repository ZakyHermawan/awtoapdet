package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type Data struct {
	Water uint `json:"water"`
	Wind  uint `json:"wind"`
}

type JSONData struct {
	Status Data `json:"status"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			var filepath = path.Join("views", "index.html")
			var tmpl, err = template.ParseFiles(filepath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jsonFile, jsonErr := os.Open("data.json")

			if jsonErr != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				err := jsonFile.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}()

			var data map[string]interface{}
			byteValue, _ := ioutil.ReadAll(jsonFile)
			err = json.Unmarshal(byteValue, &data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmp := data["status"]
			water := int(tmp.(map[string]interface{})["water"].(float64))
			wind := int(tmp.(map[string]interface{})["wind"].(float64))

			var waterStatus string
			var windStatus string

			if water <= 5 {
				waterStatus = "aman"
			} else if water < 9 {
				waterStatus = "siaga"
			} else {
				waterStatus = "bahaya"
			}

			if wind <= 6 {
				windStatus = "aman"
			} else if wind < 16 {
				windStatus = "siaga"
			} else {
				windStatus = "bahaya"
			}

			dah := struct {
				Water       int
				Wind        int
				WaterStatus string
				WindStatus  string
			}{water, wind, waterStatus, windStatus}

			err = tmpl.Execute(w, dah)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			var t Data
			err = json.Unmarshal(body, &t)
			if err != nil {
				panic(err)
			}

			jsonData := JSONData{
				Status: Data{
					Water: t.Water,
					Wind:  t.Wind,
				},
			}
			b, err := json.Marshal(jsonData)

			permissions := 0644
			err = os.WriteFile("data.json", b, os.FileMode(permissions))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var waterStatus string
			var windStatus string
			water := jsonData.Status.Water
			wind := jsonData.Status.Wind
			if water <= 5 {
				waterStatus = "aman"
			} else if water < 9 {
				waterStatus = "siaga"
			} else {
				waterStatus = "bahaya"
			}

			if wind <= 6 {
				windStatus = "aman"
			} else if wind < 16 {
				windStatus = "siaga"
			} else {
				windStatus = "bahaya"
			}

			fmt.Println("Request masuk!")
			fmt.Println("Wind:", jsonData.Status.Wind, "status:", windStatus)
			fmt.Println("Water:", jsonData.Status.Water, "status:", waterStatus)
		}
	})

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
