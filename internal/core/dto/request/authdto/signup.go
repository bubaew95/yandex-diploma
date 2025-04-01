package authdto

import "strings"

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s SignUpRequest) Validate() map[string]string {
	errorMessages := make(map[string]string)

	if strings.TrimSpace(s.Login) == "" {
		errorMessages["login_required"] = "Login is required"
	}

	if strings.TrimSpace(s.Password) == "" {
		errorMessages["password_required"] = "Password is required"
	}

	return errorMessages
}
