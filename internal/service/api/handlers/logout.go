package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/service/api/helpers"
	"gitlab.com/distributed_lab/acs/auth/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewLogoutRequest(r)
	if err != nil {
		Log(r).WithError(err).Errorf("wrong request")
		helpers.ClearTokensCookies(w)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	refreshToken, err := RefreshTokensQ(r).FilterByTokens(request.RefreshToken).Get()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to get refresh token")
		helpers.ClearTokensCookies(w)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if refreshToken == nil {
		Log(r).Errorf("no token was found in db")
		helpers.ClearTokensCookies(w)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	err = helpers.CheckValidityAndOwnerForRefreshToken(refreshToken.Token, refreshToken.OwnerId, JwtParams(r).Secret)
	if err != nil {
		Log(r).WithError(err).Errorf("something wrong with refresh token")
		helpers.ClearTokensCookies(w)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = RefreshTokensQ(r).FilterByTokens(refreshToken.Token).Delete()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to delete old refresh token")
		helpers.ClearTokensCookies(w)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	helpers.ClearTokensCookies(w)
	ape.Render(w, http.StatusOK)
}
