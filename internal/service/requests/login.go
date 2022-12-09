package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/Auth/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type LoginRequest struct {
	Data resources.Login
}

func NewLoginRequest(r *http.Request) (LoginRequest, error) {
	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request.Data); err != nil {
		return request, errors.Wrap(err, " failed to unmarshal")
	}

	return request, request.validate()
}

func (r *LoginRequest) validate() error {
	return mergeErrors(validation.Errors{
		"attributes": validation.Validate(&r.Data.Attributes, validation.Required),
	}).Filter()
}

func mergeErrors(validationErrors ...validation.Errors) validation.Errors {
	result := make(validation.Errors)
	for _, errs := range validationErrors {
		for key, err := range errs {
			result[key] = err
		}
	}
	return result
}
