package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"stori-challenge/internal/account"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var accInfo account.AccountInfoRepository

type SetAccountInfoParamsRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	godotenv.Load()

	repoType := os.Getenv("REPOSITORY_TYPE")

	repository, err := account.NewAccountInfoRepository(repoType)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	accInfo = repository

	router := mux.NewRouter()
	router.HandleFunc("/create-account", createAccount).Methods("POST")
	router.HandleFunc("/find/{id}", findAccountInfo).Methods("GET")
	router.HandleFunc("/list", listAccounts).Methods("GET")

	port := os.Getenv("ACCOUNT_INFO_PORT")

	if port == "" {
		port = "8081"
	}
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	var p SetAccountInfoParamsRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	aci := account.AccountInfo{
		Name:  p.Name,
		Email: p.Email,
	}

	result, err := accInfo.SaveAccountInfo(aci)

	if err != nil {
		http.Error(w, "error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)
}

func findAccountInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := accInfo.FindAccountInfo(id)

	if err != nil {
		http.Error(w, "error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func listAccounts(w http.ResponseWriter, _ *http.Request) {
	result, err := accInfo.List()

	if err != nil {
		http.Error(w, "error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
