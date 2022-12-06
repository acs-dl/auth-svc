package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mhrynenko/jwt_service/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type RefreshRequest struct {
	Data resources.Refresh
}

func NewRefreshRequest(r *http.Request) (RefreshRequest, error) {
	var request RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&request.Data); err != nil {
		return request, errors.Wrap(err, " failed to unmarshal")
	}

	return request, request.validate()
}

func (r *RefreshRequest) validate() error {
	return mergeErrors(validation.Errors{
		"attributes": validation.Validate(&r.Data.Attributes, validation.Required),
	}).Filter()
}
