package handlers

import (
	"errors"
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/internal/service/models"
	"gitlab.com/distributed_lab/acs/auth/internal/service/requests"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func AddPermissionUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddPermissionUserRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	exist, err := PermissionUsersQ(r).FilterByUserId(request.Data.Attributes.UserId).
		FilterByPermissionId(request.Data.Attributes.PermissionId).Get()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to get user permission")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if exist != nil {
		Log(r).Error("such user permission exists")
		ape.RenderErr(w, problems.BadRequest(errors.New("such user permission exists"))...)
		return
	}

	newPermission := data.PermissionUser{
		UserId:       request.Data.Attributes.UserId,
		PermissionId: request.Data.Attributes.PermissionId,
	}

	created, err := PermissionUsersQ(r).Create(newPermission)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to add permission for user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.PermissionUserResponse{
		Data: models.NewPermissionUserModel(created),
	}

	ape.Render(w, result)
}
