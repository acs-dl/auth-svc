package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mhrynenko/jwt_service/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type LogoutRequest struct {
	Data resources.Refresh
}

func NewLogoutRequest(r *http.Request) (LogoutRequest, error) {
	var request LogoutRequest

	if err := json.NewDecoder(r.Body).Decode(&request.Data); err != nil {
		return request, errors.Wrap(err, " failed to unmarshal")
	}

	return request, request.validate()
}

func (r *LogoutRequest) validate() error {
	return mergeErrors(validation.Errors{
		"attributes": validation.Validate(&r.Data.Attributes, validation.Required),
	}).Filter()
}
