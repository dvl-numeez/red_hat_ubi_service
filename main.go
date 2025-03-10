package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var redHatServiceDb *RedHatDBService
var PORT string

func main() {

	if PORT = os.Getenv("PORT"); PORT == "" {
		log.Fatal("PORT env required")
	}
	if MONGO_URI = os.Getenv("MONGO_URI"); MONGO_URI == "" {
		log.Fatal("MONGO_URI env variable required")
	}
	if COLLECTION_NAME = os.Getenv("COLLECTION_NAME"); COLLECTION_NAME == "" {
		log.Fatal("COLLECTION_NAME env variable required")
	}
	if DB_NAME = os.Getenv("DB_NAME"); DB_NAME == "" {
		log.Fatal("DB_NAME env variable required")
	}

	redHatServiceDb = NewRedHatDBService(COLLECTION_NAME)
	router := http.NewServeMux()
	router.HandleFunc("/health", health)
	router.HandleFunc("/users", handlerGetUsers)
	router.HandleFunc("/addUser", handlerAddUser)
	log.Println("Server running on the port: ", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), router); err != nil {
		log.Fatal("Unable to run server due to err: ", err)
	}

}

func health(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Service running",
	}
	respByte, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("Server crashed due to error: ", err)
	}
	w.Header().Add("content-type", "json/application")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respByte)

}

func handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := redHatServiceDb.fetchUsers()
	if err != nil {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "Unable to fetch users",
		}
		respByte, err := json.Marshal(resp)
		if err != nil {
			log.Fatal("Server crashed due to error: ", err)
		}
		w.Header().Add("content-type", "json/application")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(respByte)
		return
	}

	data, err := json.Marshal(users)
	if err != nil {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "Unable to fetch users",
		}
		respByte, err := json.Marshal(resp)
		if err != nil {
			log.Fatal("Server crashed due to error: ", err)
		}
		w.Header().Add("content-type", "json/application")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(respByte)
		return
	}
	w.Header().Add("content-type", "json/application")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)

}

func handlerAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "Invalid method",
		}
		respByte, err := json.Marshal(resp)
		if err != nil {
			log.Fatal("Server crashed due to error: ", err)
		}
		w.Header().Add("content-type", "json/application")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(respByte)
		return
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		respByte, err := json.Marshal(resp)
		if err != nil {
			log.Fatal("Server crashed due to error: ", err)
		}
		w.Header().Add("content-type", "json/application")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(respByte)
		return
	}
	if err := redHatServiceDb.addUser(user); err != nil {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		respByte, err := json.Marshal(resp)
		if err != nil {
			log.Fatal("Server crashed due to error: ", err)
		}
		w.Header().Add("content-type", "json/application")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(respByte)
		return
	}
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "User added successfully",
	}
	respByte, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("Server crashed due to error: ", err)
	}
	w.Header().Add("content-type", "json/application")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respByte)

}
