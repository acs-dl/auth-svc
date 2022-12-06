package handlers

import (
	"github.com/mhrynenko/jwt_service/internal/service/helpers"
	"github.com/mhrynenko/jwt_service/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewLogoutRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	refreshToken, err := RefreshTokensQ(r).FilterByToken(request.Data.Attributes.Token).Get()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to get refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = helpers.CheckRefreshToken(refreshToken.Token, refreshToken.OwnerId)
	if err != nil {
		Log(r).WithError(err).Info("something wrong with refresh token")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = RefreshTokensQ(r).Delete(refreshToken.Token)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to delete old refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK)
}
