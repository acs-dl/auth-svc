package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/auth/internal/service/requests"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewLoginRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := UsersQ(r).FilterByEmail(request.Data.Attributes.Email).Get()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if user == nil {
		Log(r).Errorf("no user with such email `%s`", request.Data.Attributes.Email)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Data.Attributes.Password)) != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	jwt := JwtParams(r)

	permissions, err := PermissionUsersQ(r).FilterByUserId(user.Id).Select()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to get user permissions")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	permissionsString, err := helpers.CreatePermissionsString(permissions, PermissionsQ(r))
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

	refresh, err, claims := helpers.GenerateRefreshToken(*user, helpers.ParseToUnix(jwt.RefreshLife), jwt.Secret, permissionsString)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to create refresh")
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
