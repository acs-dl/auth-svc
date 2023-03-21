package requests

import (
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"net/http"
)

type DeletePermissionRequest struct {
	ModuleName     string `json:"module_name"`
	PermissionName string `json:"permission_name"`
}

func NewDeletePermissionRequest(r *http.Request) (DeletePermissionRequest, error) {
	var request DeletePermissionRequest

	request.ModuleName = chi.URLParam(r, "name")
	request.PermissionName = chi.URLParam(r, "permission_name")

	return request, request.validate()
}

func (r *DeletePermissionRequest) validate() error {
	return validation.Errors{
		"module_name":     validation.Validate(&r.ModuleName, validation.Required),
		"permission_name": validation.Validate(&r.PermissionName, validation.Required),
	}.Filter()
}
