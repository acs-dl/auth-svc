package requests

import (
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"net/http"
)

type DeleteModuleRequest struct {
	ModuleName string `json:"module_name"`
}

func NewDeleteModuleRequest(r *http.Request) (DeleteModuleRequest, error) {
	var request DeleteModuleRequest

	request.ModuleName = chi.URLParam(r, "name")

	return request, request.validate()
}

func (r *DeleteModuleRequest) validate() error {
	return validation.Errors{
		"name": validation.Validate(&r.ModuleName, validation.Required),
	}.Filter()
}
