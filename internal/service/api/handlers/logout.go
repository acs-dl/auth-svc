package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/service/api/helpers"
	"gitlab.com/distributed_lab/acs/auth/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewLogoutRequest(r)
	if err != nil {
		Log(r).WithError(err).Errorf("wrong request")
		helpers.ClearTokensCookies(w)
		ape.Render(w, http.StatusOK)
		return
	}

	refreshToken, err := RefreshTokensQ(r).FilterByTokens(request.RefreshToken).Get()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to get refresh token")
		helpers.ClearTokensCookies(w)
		ape.Render(w, http.StatusOK)
		return
	}
	if refreshToken == nil {
		Log(r).Errorf("no token was found in db")
		helpers.ClearTokensCookies(w)
		ape.Render(w, http.StatusOK)
		return
	}

	err = helpers.CheckValidityAndOwnerForRefreshToken(refreshToken.Token, refreshToken.OwnerId, JwtParams(r).Secret)
	if err != nil {
		Log(r).WithError(err).Errorf("something wrong with refresh token")
		helpers.ClearTokensCookies(w)
		ape.Render(w, http.StatusOK)
		return
	}

	err = RefreshTokensQ(r).FilterByTokens(refreshToken.Token).Delete()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to delete old refresh token")
		helpers.ClearTokensCookies(w)
		ape.Render(w, http.StatusOK)
		return
	}

	helpers.ClearTokensCookies(w)
	ape.Render(w, http.StatusOK)
}
