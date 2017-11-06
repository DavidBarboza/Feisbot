package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Port    string `json:"port"`
	CertPem string `json:"cert_pem"`
	KeyPem  string `json:"key_pem"`
	MyToken string `json:"my_token"`
}

var config Config

func main() {
	loadConfig()

	http.HandleFunc("/", saludar)
	http.HandleFunc("/fbwebhook",fbwebhook);

	log.Printf("servidor iniciado en https://localhost%s", config.Port)
	err := http.ListenAndServeTLS(config.Port,
					config.CertPem,
					config.KeyPem,
					nil)
	if err != nil {
		log.Println(err)
	}
	//log.Println(http.ListenAndServe(":8080", nil))
}

func saludar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func loadConfig() {
	log.Println("Loading configuration file...")
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("Error converting config file: %v", err)
	}
	log.Println("Configuration loaded")
}

func fbwebhook(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet{
		vt := r.URL.Query().Get("hub.verify_token")
		if vt == config.MyToken{
			hc := r.URL.Query().Get("hub.challenge")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(hc))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Token is not valid"))
		return
	}
}
