package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AddModuleUserRequest struct {
	Data resources.ModuleUser
}

func NewAddModuleUserRequest(r *http.Request) (AddModuleUserRequest, error) {
	var request AddModuleUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, " failed to unmarshal")
	}

	return request, request.validate()
}

func (r *AddModuleUserRequest) validate() error {
	return mergeErrors(validation.Errors{
		"attributes": validation.Validate(&r.Data.Attributes, validation.Required),
	}).Filter()
}
