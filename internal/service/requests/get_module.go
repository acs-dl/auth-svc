package requests

import (
	"github.com/go-chi/chi"
	"net/http"
)

type GetModuleRequest struct {
	ModuleName string `json:"module_name"`
}

func NewGetModuleRequestRequest(r *http.Request) (GetModuleRequest, error) {
	var request GetModuleRequest

	request.ModuleName = chi.URLParam(r, "name")

	return request, nil
}
