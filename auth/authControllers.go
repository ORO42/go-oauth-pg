package auth

import (
	"fmt"
	"gop/users"
	"net/http"

	"github.com/markbates/goth/gothic"
)

// Start OAuth flow
func BeginAuthController(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

// Finish OAuth flow and store user info in session
func CompleteAuthController(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// Store user info in session
	session, err := gothic.Store.Get(r, "session-store")
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	session.Values["auth_provider_title"] = user.Provider
	session.Values["auth_provider_user_id"] = user.UserID
	session.Values["auth_provider_user_name"] = user.Name
	session.Values["auth_provider_email"] = user.Email
	session.Values["auth_provider_avatar_url"] = user.AvatarURL

	// Add user to db if doesn't already exist
	DBUserID, err := users.CreateDBUser(user.Provider, user.UserID, user.Email)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// Add the DB user to the session
	session.Values["db_user_id"] = DBUserID

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	fmt.Printf("User session saved: %+v\n", user)

	// Redirect to dashboard or return user info
	http.Redirect(w, r, "/test", http.StatusFound)
}

// Clear session
func LogoutController(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Logged out successfully")
}
