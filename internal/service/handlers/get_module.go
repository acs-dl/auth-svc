package handlers

import (
	"gitlab.com/distributed_lab/acs/auth/internal/service/models"
	"gitlab.com/distributed_lab/acs/auth/internal/service/requests"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetModule(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetModuleRequestRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := PermissionsQ(r).WithModules().FilterByModuleName(request.ModuleName).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get module")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := resources.ModulePermissionListResponse{
		Data: models.NewModulePermissionsList(result),
	}

	ape.Render(w, response)
}
