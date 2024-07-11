package handlers

import (
	"encoding/json"
	"net/http"
	"vokki_cloud/internal/models"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser models.NewUser
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, _ := models.GetUser(newUser.Email)

	if user.Email != "" {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

}
