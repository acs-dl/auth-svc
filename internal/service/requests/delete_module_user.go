package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type DeletePermissionUserRequest struct {
	Data resources.PermissionUser
}

func NewDeletePermissionUserRequest(r *http.Request) (DeletePermissionUserRequest, error) {
	var request DeletePermissionUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, " failed to unmarshal")
	}

	return request, request.validate()
}

func (r *DeletePermissionUserRequest) validate() error {
	return validation.Errors{
		"permission_id": validation.Validate(&r.Data.Attributes.PermissionId, validation.Required),
		"user_id":       validation.Validate(&r.Data.Attributes.UserId, validation.Required),
	}.Filter()
}
