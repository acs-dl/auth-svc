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

func AddModuleUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddModuleUserRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	newPermission := data.ModuleUser{
		UserId:   request.Data.Attributes.UserId,
		ModuleId: request.Data.Attributes.ModuleId,
	}

	created, err := ModulesUsersQ(r).Create(newPermission)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to add permission for user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.ModuleUserResponse{
		Data: models.NewModuleUserModel(created),
	}

	ape.Render(w, result)
}
