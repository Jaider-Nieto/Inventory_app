package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidationMiddleware(next http.Handler, model interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Restore the body so the next handler can read it
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Unmarshal the body into the model
		if err := json.Unmarshal(bodyBytes, model); err != nil {
			http.Error(w, "Error unmarshalling request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(model); err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Field validation error on " + err.Field() + ": " + err.Tag()))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
