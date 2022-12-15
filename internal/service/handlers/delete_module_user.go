package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteModuleUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewDeleteModuleUserRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	moduleUser := data.ModuleUser{
		UserId:   request.Data.Attributes.UserId,
		ModuleId: request.Data.Attributes.ModuleId,
	}

	err = ModulesUsersQ(r).Delete(moduleUser)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to remove permission for user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK)
}
