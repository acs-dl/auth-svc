package requests

import (
	"github.com/go-chi/chi"
	"net/http"
)

type DeleteModuleRequest struct {
	ModuleName string `json:"module_name"`
}

func NewDeleteModuleRequest(r *http.Request) (DeleteModuleRequest, error) {
	var request DeleteModuleRequest

	request.ModuleName = chi.URLParam(r, "name")

	return request, nil
}
