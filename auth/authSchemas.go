package auth

type SessionUser struct {
	AuthProviderTitle    string `json:"auth_provider_title"`
	AuthProviderUserID   string `json:"auth_provider_user_id"`
	AuthProviderUserName string `json:"auth_provider_user_name"`
	AuthProviderEmail    string `json:"auth_provider_email"`
	DBUserID             int    `json:"db_user_id"`
}
