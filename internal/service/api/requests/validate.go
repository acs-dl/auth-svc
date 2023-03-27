package requests

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
)

type ValidateRequest struct {
	Token string `json:"token"`
}

func NewValidateRequest(r *http.Request) (ValidateRequest, error) {
	var request ValidateRequest

	splitedAuthHeader := strings.Split(r.Header.Get("Authorization"), " ")
	if len(splitedAuthHeader) != 2 {
		return request, errors.New("no token was provided")
	}

	request.Token = splitedAuthHeader[1]

	return request, request.validate()
}

func (r *ValidateRequest) validate() error {
	return validation.Errors{
		"token": validation.Validate(&r.Token, validation.Required, validation.Match(regexp.MustCompile(data.TokenRegExpStr))),
	}.Filter()
}
