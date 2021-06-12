package middlewares

import (
	"encoding/json"
	"net/http"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	raw, _ := json.Marshal(data)
	if statusCode == http.StatusFound {
		w.Header().Add("Location", "https://go.go")
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(raw)
}

func JsonRedirectResponse(w http.ResponseWriter, url string) {
	raw, _ := json.Marshal(url)
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusFound)
	_, _ = w.Write(raw)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, err error, message string) {
	if err != nil {
		JsonResponse(w, statusCode, struct {
			Error string `json:"error"`
			Message string `json:"message"`
		}{
			Error: err.Error(),
			Message: message,
		})
	} else {
		JsonResponse(w, statusCode, struct {
			Message string `json:"message"`
		}{
			Message: message,
		})
	}
}