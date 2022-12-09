package handlers

import (
	"gitlab.com/distributed_lab/Auth/internal/data"
	"gitlab.com/distributed_lab/Auth/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
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
		Log(r).WithError(err).Error(err, "failed to add module")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK)
}
