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
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Request body is empty"))
            return
        }

        bodyBytes, err := io.ReadAll(r.Body)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Error reading request body: " + err.Error()))
            return
        }

        r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

        if err := json.Unmarshal(bodyBytes, model); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Error unmarshalling request body: " + err.Error()))
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
