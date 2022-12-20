package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/auth/internal/service/requests"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
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
	if refreshToken == nil {
		Log(r).Info("no token was found in db")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	jwt := JwtParams(r)

	err = helpers.CheckRefreshToken(refreshToken.Token, refreshToken.OwnerId, jwt.Secret)
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

	permissions, err := ModulesUsersQ(r).FilterByUserId(user.Id).Select()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to get user permissions")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	permissionsString, err := helpers.CreatePermissionsString(permissions, ModulesQ(r))
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to get create user permissions string")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	access, err := helpers.GenerateAccessToken(*user, helpers.ParseToUnix(jwt.AccessLife), jwt.Secret, permissionsString)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to create access token")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	err = AmountsQ(r).Add("access")
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to add counter user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	refresh, err, claims := helpers.GenerateRefreshToken(*user, helpers.ParseToUnix(jwt.RefreshLife), jwt.Secret, permissionsString)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to create refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	err = AmountsQ(r).Add("refresh")
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to add counter")
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

	result := resources.AuthTokenResponse{
		Data: resources.AuthToken{
			Attributes: resources.AuthTokenAttributes{
				Access:  access,
				Refresh: refresh,
			},
		},
	}

	ape.Render(w, result)
}
