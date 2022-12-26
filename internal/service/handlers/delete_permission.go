package handlers

import (
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func DeletePermission(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewDeletePermissionRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	module, err := ModulesQ(r).GetByName(request.ModuleName)
	if err != nil {
		Log(r).WithError(err).Error("failed to get module")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if module == nil {
		Log(r).Error("no such module")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	err = PermissionsQ(r).Delete(data.Permission{
		ModuleId: module.Id,
		Name:     request.PermissionName,
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to delete permission")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusAccepted)
}
