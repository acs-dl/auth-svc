package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteModule(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewDeleteModuleRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = ModulesQ(r).Delete(request.ModuleName)
	if err != nil {
		Log(r).WithError(err).Error("failed to delete module")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusAccepted)
}
