package handlers

import "net/http"

func LandingPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/views/index.html")

}
