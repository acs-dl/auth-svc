package requests

import (
	"github.com/go-chi/chi"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

type GetModuleRequest struct {
	ModuleName string `json:"module_name"`
}

func NewGetModuleRequestRequest(r *http.Request) (GetModuleRequest, error) {
	var request GetModuleRequest

	request.ModuleName = chi.URLParam(r, "name")

	return request, nil
}

func (r *GetModuleRequest) validate() error {
	return validation.Errors{
		"name": validation.Validate(&r.ModuleName, validation.Required),
	}.Filter()
}
