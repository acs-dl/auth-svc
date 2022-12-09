package handlers

import (
	"gitlab.com/distributed_lab/Auth/internal/data"
	"gitlab.com/distributed_lab/Auth/internal/service/models"
	"gitlab.com/distributed_lab/Auth/internal/service/requests"
	"gitlab.com/distributed_lab/Auth/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
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
