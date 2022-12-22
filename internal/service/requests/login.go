package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/acs/auth/resources"
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
		"attributes": validation.Validate(&r.Data.Attributes, validation.Required),
	}.Filter()
}
