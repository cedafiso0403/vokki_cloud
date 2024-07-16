package handlers

import (
	"encoding/json"
	"net/http"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/services"
)

func Login(w http.ResponseWriter, r *http.Request) {

	credentials := services.Credentials{}

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	err := decoder.Decode(&credentials)

	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if credentials.Email == "" || credentials.Password == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	userID, token, err := services.Authenticate(credentials)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = models.StoreToken(int64(userID), token, vokki_constants.AuthToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
