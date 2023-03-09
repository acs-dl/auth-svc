package handlers

import (
	"gitlab.com/distributed_lab/acs/auth/internal/service/models"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/internal/service/requests"
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

	newModule := data.Module{
		Name: request.Data.Attributes.Module,
	}

	module, err := ModulesQ(r).Create(newModule)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to create module")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if module == nil {
		Log(r).WithError(err).Error(err, "no such module")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	newPermission := data.Permission{
		ModuleId: module.Id,
		Name:     request.Data.Attributes.Permission,
	}

	permission, err := PermissionsQ(r).Create(newPermission)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to create permission")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if permission == nil {
		Log(r).WithError(err).Error(err, "no such permission")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	result := resources.ModulePermissionResponse{
		Data: models.NewModulePermissionModel(*module, *permission),
	}

	ape.Render(w, result)
}
