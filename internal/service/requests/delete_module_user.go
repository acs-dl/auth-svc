package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/Auth/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type DeleteModuleUserRequest struct {
	Data resources.ModuleUser
}

func NewDeleteModuleUserRequest(r *http.Request) (DeleteModuleUserRequest, error) {
	var request DeleteModuleUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request.Data); err != nil {
		return request, errors.Wrap(err, " failed to unmarshal")
	}

	return request, request.validate()
}

func (r *DeleteModuleUserRequest) validate() error {
	return mergeErrors(validation.Errors{
		"attributes": validation.Validate(&r.Data.Attributes, validation.Required),
	}).Filter()
}
