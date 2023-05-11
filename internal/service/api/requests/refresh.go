package requests

import (
	"net/http"
	"regexp"

	"github.com/acs-dl/auth-svc/internal/data"
	validation "github.com/go-ozzo/ozzo-validation"
)

type RefreshRequest struct {
	RefreshToken string
}

func NewRefreshRequest(r *http.Request) (RefreshRequest, error) {
	var request RefreshRequest

	cookie, err := r.Cookie(data.RefreshCookie)
	if err != nil {
		return request, err
	}
	request.RefreshToken = cookie.Value

	return request, request.validate()
}

func (r *RefreshRequest) validate() error {
	return validation.Errors{
		"token": validation.Validate(&r.RefreshToken, validation.Required, validation.Match(regexp.MustCompile(data.TokenRegExpStr))),
	}.Filter()
}
