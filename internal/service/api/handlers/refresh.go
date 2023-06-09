package handlers

import (
	"net/http"

	"github.com/acs-dl/auth-svc/internal/service/api/helpers"
	"github.com/acs-dl/auth-svc/internal/service/api/models"
	"github.com/acs-dl/auth-svc/internal/service/api/requests"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/acs-dl/auth-svc/internal/data"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewRefreshRequest(r)
	if err != nil {
		Log(r).WithError(err).Errorf("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	refreshToken, err := checkRefreshToken(RefreshTokensQ(r), request.RefreshToken, JwtParams(r).Secret)
	if err != nil {
		Log(r).WithError(err).Error(err, " failed to check refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	user, err := UsersQ(r).FilterByIds(refreshToken.OwnerId).Get()
	if err != nil {
		Log(r).WithError(err).Error(err, " failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user == nil {
		Log(r).Error("no such user")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	permissionsString, err := getPermissionsString(PermissionsQ(r), user.Status)
	if err != nil {
		Log(r).WithError(err).Error(err, " failed to get permissions string")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	access, refresh, claims, err := generateTokens(data.GenerateTokens{
		User:              *user,
		AccessLife:        helpers.ParseDurationStringToUnix(JwtParams(r).AccessLife),
		RefreshLife:       helpers.ParseDurationStringToUnix(JwtParams(r).RefreshLife),
		Secret:            JwtParams(r).Secret,
		PermissionsString: permissionsString,
	})
	if err != nil {
		Log(r).WithError(err).Errorf(" failed to generate access and refresh tokens")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = RefreshTokensQ(r).FilterByTokens(refreshToken.Token).Delete()
	if err != nil {
		Log(r).WithError(err).Error(err, " failed to delete old refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	newRefreshToken := data.RefreshToken{
		Token:     refresh,
		OwnerId:   claims.OwnerId,
		ValidTill: claims.ExpiresAt,
	}

	err = RefreshTokensQ(r).Create(newRefreshToken)
	if err != nil {
		Log(r).WithError(err).Error(err, " failed to create instance of refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	helpers.SetTokensCookies(w, access, refresh)
	ape.Render(w, models.NewAuthTokenResponse(access))
}

func checkRefreshToken(refreshTokensQ data.RefreshTokens, token, secret string) (*data.RefreshToken, error) {
	refreshToken, err := refreshTokensQ.FilterByTokens(token).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get refresh token")
	}

	if refreshToken == nil {
		return nil, errors.Errorf("no token was found in db")
	}

	err = helpers.CheckValidityAndOwnerForRefreshToken(refreshToken.Token, refreshToken.OwnerId, secret)
	if err != nil {
		return nil, errors.Wrap(err, "something wrong with refresh token")
	}

	return refreshToken, nil
}
