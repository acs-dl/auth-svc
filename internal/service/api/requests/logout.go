package requests

import (
	"net/http"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
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
