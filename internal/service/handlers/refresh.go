package handlers

import (
	"github.com/mhrynenko/jwt_service/internal/data"
	"github.com/mhrynenko/jwt_service/internal/service/helpers"
	"github.com/mhrynenko/jwt_service/internal/service/requests"
	"github.com/mhrynenko/jwt_service/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewRefreshRequest(r)
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

	user, err := UsersQ(r).FilterById(refreshToken.OwnerId).Get()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	refresh, err, claims := helpers.GenerateRefreshToken(*user)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to create refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = RefreshTokensQ(r).Delete(refreshToken.Token)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to delete old refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	newRefreshToken := data.RefreshToken{
		Token:     refresh,
		OwnerId:   claims["owner_id"].(int64),
		ValidDate: claims["exp"].(int64),
	}

	err = RefreshTokensQ(r).Create(newRefreshToken)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to create instance of refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.RefreshResponse{
		Data: resources.Refresh{
			Attributes: resources.RefreshAttributes{
				Token: refresh,
			},
		},
	}

	ape.Render(w, result)
}
