package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pusher/pusher-http-go/v5"
)

var pusherClient pusher.Client

func main() {
	pusherClient = pusher.Client{
		AppID:   "1768708",
		Key:     "7fddf3fea8628e69ba6b",
		Secret:  "fa91f8bf7cd0893e705e",
		Cluster: "ap1",
		Secure:  true,
	}

	http.HandleFunc("/send-message", sendMessageHandler)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Set header untuk izinkan CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Pastikan hanya menerima metode POST atau OPTIONS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Dekode JSON dari body request
	var payload struct {
		Event   string            `json:"event"`
		Message map[string]string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding JSON: %v", err)
		return
	}

	// Kirim pesan ke Pusher
	err := pusherClient.Trigger("progsil-web", payload.Event, payload.Message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error sending message: %v", err)
		return
	}

	// Berhasil mengirim pesan
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message sent successfully")
}
