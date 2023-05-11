package requests

import (
	"net/http"
	"regexp"

	"github.com/acs-dl/auth-svc/internal/data"
	validation "github.com/go-ozzo/ozzo-validation"
)

type LogoutRequest struct {
	RefreshToken string
}

func NewLogoutRequest(r *http.Request) (LogoutRequest, error) {
	var request LogoutRequest
	cookie, err := r.Cookie(data.RefreshCookie)
	if err != nil {
		return request, err
	}
	request.RefreshToken = cookie.Value

	return request, request.validate()
}

func (r *LogoutRequest) validate() error {
	return validation.Errors{
		"token": validation.Validate(&r.RefreshToken, validation.Required, validation.Match(regexp.MustCompile(data.TokenRegExpStr))),
	}.Filter()
}
