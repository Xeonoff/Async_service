package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	//"strings"
	"time"
)

const BackendURL = "http://localhost:8000/"
const Password = "abrakadabra"

type Request struct {
	ID string `json:"id_sending"`
}

type Response struct {
	ID       string `json:"id_sending"`
	Status   string `json:"rand_status"`
	Password string `json:"password"`
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statuses := []string{"P", "D", "A", "W"}
	randomStatus := statuses[rand.Intn(len(statuses))]

	response := Response{
		ID:       req.ID,
		Status:   randomStatus,
		Password: Password,
	}
	respJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Request received successfully") // Отправляем ответ немедленно

	go func() {
		time.Sleep(5 * time.Second) // Ждем 5 секунд перед отправкой PUT запроса

		client := &http.Client{}
		request, err := http.NewRequest("PUT", BackendURL+"send/", bytes.NewBuffer(respJSON))
		if err != nil {
			log.Println(err.Error()) // Логируем ошибку вместо отправки в ответ, так как ответ уже отправлен
			return
		}
		request.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(request)
		if err != nil {
			log.Println(err.Error()) // Логируем ошибку
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Println("Request processed successfully")
		} else {
			log.Printf("Backend returned error: %vn", resp.Status)
		}
	}()
}

func main() {
	http.HandleFunc("/send/", func(w http.ResponseWriter, r *http.Request) {
		postHandler(w, r)
	})

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
