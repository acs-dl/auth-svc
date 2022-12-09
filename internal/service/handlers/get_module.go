package handlers

import (
	"gitlab.com/distributed_lab/Auth/internal/service/models"
	"gitlab.com/distributed_lab/Auth/internal/service/requests"
	"gitlab.com/distributed_lab/Auth/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetModule(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetModuleRequestRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	modulePermissions, err := ModulesQ(r).FilterByModule(request.ModuleName).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get module")
		ape.Render(w, problems.InternalError())
		return
	}

	response := resources.ModulePermissionListResponse{
		Data: models.NewModulePermissionsList(modulePermissions),
	}

	ape.Render(w, response)
}
