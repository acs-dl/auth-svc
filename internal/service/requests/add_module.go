package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AddModuleRequest struct {
	Data resources.ModulePermission
}

func NewAddModuleRequest(r *http.Request) (AddModuleRequest, error) {
	var request AddModuleRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, " failed to unmarshal")
	}

	return request, request.validate()
}

func (r *AddModuleRequest) validate() error {
	return mergeErrors(validation.Errors{
		"attributes": validation.Validate(&r.Data.Attributes, validation.Required),
	}).Filter()
}
