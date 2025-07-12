package auth

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func InitAuth() {
	key := os.Getenv("SESSION_SECRET")
	maxAge := 86400 * 30 // 30 days
	isProd := false      // set to true in production

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	// Configure the Google OAuth provider
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
	)

	// Required if you're not using gorilla/mux or pat
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		return "google", nil
	}
}

// Helper function to get user from session
func GetUserFromSession(r *http.Request) (SessionUser, error) {
	session, err := gothic.Store.Get(r, "session-store")
	if err != nil {
		return SessionUser{}, err
	}

	user := SessionUser{}

	if v, ok := session.Values["auth_provider_title"].(string); ok {
		user.AuthProviderTitle = v
	}
	if v, ok := session.Values["auth_provider_user_id"].(string); ok {
		user.AuthProviderUserID = v
	}
	if v, ok := session.Values["auth_provider_user_name"].(string); ok {
		user.AuthProviderUserName = v
	}
	if v, ok := session.Values["auth_provider_email"].(string); ok {
		user.AuthProviderEmail = v
	}
	if v, ok := session.Values["db_user_id"].(int); ok {
		user.DBUserID = v
	}

	return user, nil
}
