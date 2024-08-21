package handlers

import "net/http"

func TermAndConditions(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/views/terms_and_conditions.html")
}
