package auth

import (
	"net/http"

	"github.com/markbates/goth/gothic"
)

// Authentication middleware
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := gothic.Store.Get(r, "session-store")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		// Check if user is authenticated
		userID, exists := session.Values["user_id"]
		if !exists || userID == nil {
			http.Error(w, "Unauthorized - Please login first", http.StatusUnauthorized)
			return
		}

		// User is authenticated, proceed to the handler
		next(w, r)
	}
}
