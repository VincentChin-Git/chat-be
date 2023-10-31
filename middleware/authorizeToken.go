package middleware

import (
	"chat-be/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type ContextKey string

const TokenKey = ContextKey("parsedId")

func authorizeToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path

		// Do stuff
		var skippedPath = []string{"/api/user/signup", "/api/user/login"}
		for _, path := range skippedPath {
			if path == reqPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			utils.JsonResponseError(w, "888888", "Invalid Token", http.StatusUnauthorized)
			return
		}

		result, err := utils.DecodeToken(authHeader[1])
		if err != nil {
			utils.JsonResponseError(w, "888888", "Invalid Token", http.StatusUnauthorized)
			return
		}

		fmt.Println("_id:", result)

		r = r.WithContext(context.WithValue(r.Context(), TokenKey, result))

		next.ServeHTTP(w, r)
	})
}
