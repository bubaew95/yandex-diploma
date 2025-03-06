package authdto

import "strings"

type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s SignInRequest) Validate() map[string]string {
	errorMessages := make(map[string]string)

	if strings.TrimSpace(s.Login) == "" {
		errorMessages["login_required"] = "Login is empty"
	}

	if strings.TrimSpace(s.Password) == "" {
		errorMessages["password_required"] = "Password is empty"
	}

	return errorMessages
}
