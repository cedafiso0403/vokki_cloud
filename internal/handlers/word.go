package handlers

import (
	"log"
	"net/http"
	"time"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/httputil"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/utils"
)

func GetWordTranslations(w http.ResponseWriter, r *http.Request) {

	timeNow := time.Now().UTC()

	userID := r.Context().Value(vokki_constants.UserIDKey)

	if userID == nil {
		httpError := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Message:   "User ID not found",
			Status:    http.StatusBadRequest,
		}
		httputil.ErrorJsonResponse(w, httpError, http.StatusBadRequest)
		return
	}

	word := r.URL.Query().Get("word")

	if word == "" {
		httpError := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Message:   "Word parameter is required",
			Status:    http.StatusBadRequest,
		}
		httputil.ErrorJsonResponse(w, httpError, http.StatusBadRequest)
		return
	}

	// if lang is not provided, return all FOR NOW AS WE DO NOT HAVE A DEFAULT LANGUAGE
	//lang := r.URL.Query().Get("lang")

	translations, err := models.GetWordTranslations(word, "en")

	if err != nil {
		log.Println("Error getting word translations: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	httputil.SuccessJsonResponse(w, map[string]any{
		"data":       translations,
		"time_stamp": utils.FormatDate(timeNow),
	})

}
