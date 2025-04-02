package models

// UserInfo represents a user profile according to OpenID Connect standards
// as returned by Google's OAuth2 userinfo endpoint
type UserInfo struct {
	// Standard OpenID Connect fields
	Sub           string `json:"sub"`            // Subject identifier - unique ID for the user
	Name          string `json:"name"`           // Full name
	GivenName     string `json:"given_name"`     // First name
	FamilyName    string `json:"family_name"`    // Last name
	Email         string `json:"email"`          // Email address
	EmailVerified bool   `json:"email_verified"` // Whether email is verified
	Picture       string `json:"picture"`        // URL to profile picture

	// Optional additional fields
	Locale string `json:"locale,omitempty"` // User's locale/language
	HD     string `json:"hd,omitempty"`     // Hosted domain (for G Suite users)
}

// NewDefaultUser creates a user with default values
func NewDefaultUser() *UserInfo {
	return &UserInfo{
		Sub:           "123456789",
		Name:          "Test User",
		GivenName:     "Test",
		FamilyName:    "User",
		Email:         "testuser@example.com",
		EmailVerified: true,
		Picture:       "https://example.com/profile.jpg",
	}
}

// Clone creates a deep copy of the user
func (u *UserInfo) Clone() *UserInfo {
	if u == nil {
		return nil
	}

	return &UserInfo{
		Sub:           u.Sub,
		Name:          u.Name,
		GivenName:     u.GivenName,
		FamilyName:    u.FamilyName,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Picture:       u.Picture,
		Locale:        u.Locale,
		HD:            u.HD,
	}
}

// Merge updates this user with non-zero values from the other user
func (u *UserInfo) Merge(other *UserInfo) {
	if other == nil {
		return
	}

	if other.Sub != "" {
		u.Sub = other.Sub
	}
	if other.Name != "" {
		u.Name = other.Name
	}
	if other.GivenName != "" {
		u.GivenName = other.GivenName
	}
	if other.FamilyName != "" {
		u.FamilyName = other.FamilyName
	}
	if other.Email != "" {
		u.Email = other.Email
	}
	// EmailVerified is boolean, so always update it
	u.EmailVerified = other.EmailVerified

	if other.Picture != "" {
		u.Picture = other.Picture
	}
	if other.Locale != "" {
		u.Locale = other.Locale
	}
	if other.HD != "" {
		u.HD = other.HD
	}
}

// UpdateFromConfig updates user info from a configuration struct
// This will be used by the /config endpoint to modify user data dynamically
func UpdateUserFromConfig(user *UserInfo, config map[string]interface{}) {
	if config == nil {
		return
	}

	// Update fields from the config map
	if sub, ok := config["sub"].(string); ok {
		user.Sub = sub
	}
	if name, ok := config["name"].(string); ok {
		user.Name = name
	}
	if givenName, ok := config["given_name"].(string); ok {
		user.GivenName = givenName
	}
	if familyName, ok := config["family_name"].(string); ok {
		user.FamilyName = familyName
	}
	if email, ok := config["email"].(string); ok {
		user.Email = email
	}
	if emailVerified, ok := config["email_verified"].(bool); ok {
		user.EmailVerified = emailVerified
	}
	if picture, ok := config["picture"].(string); ok {
		user.Picture = picture
	}
	if locale, ok := config["locale"].(string); ok {
		user.Locale = locale
	}
	if hd, ok := config["hd"].(string); ok {
		user.HD = hd
	}
}
