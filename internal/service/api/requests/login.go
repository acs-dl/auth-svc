package requests

import (
	"encoding/json"
	"net/http"

	"github.com/acs-dl/auth-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type LoginRequest struct {
	Data resources.Login
}

func NewLoginRequest(r *http.Request) (LoginRequest, error) {
	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, " failed to unmarshal")
	}

	return request, request.validate()
}

func (r *LoginRequest) validate() error {
	return validation.Errors{
		"email":    validation.Validate(&r.Data.Attributes.Email, validation.Required, validation.Length(1, 255)),
		"password": validation.Validate(&r.Data.Attributes.Password, validation.Required),
	}.Filter()
}
