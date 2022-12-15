package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/internal/service/models"
	"gitlab.com/distributed_lab/acs/auth/internal/service/requests"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func AddModule(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddModuleRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	newModule := data.ModulePermission{
		Permission: request.Data.Attributes.Permission,
		ModuleName: request.Data.Attributes.Module,
	}

	createdModule, err := ModulesQ(r).Create(newModule)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to add module")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.ModulePermissionResponse{
		Data: models.NewModulePermissionModel(createdModule),
	}

	ape.Render(w, result)
}
