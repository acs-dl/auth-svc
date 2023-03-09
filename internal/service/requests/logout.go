package requests

import (
	"encoding/json"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"net/http"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type LogoutRequest struct {
	Data resources.Refresh
}

func NewLogoutRequest(r *http.Request) (LogoutRequest, error) {
	var request LogoutRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, " failed to unmarshal")
	}

	return request, request.validate()
}

func (r *LogoutRequest) validate() error {
	return validation.Errors{
		"token": validation.Validate(&r.Data.Attributes.Token, validation.Required, validation.Match(regexp.MustCompile(data.TokenRegExpStr))),
	}.Filter()
}
