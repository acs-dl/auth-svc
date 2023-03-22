package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/service/models"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/auth/internal/service/requests"
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

	user, err := checkUserAndPassword(request, UsersQ(r))
	if err != nil {
		Log(r).WithError(err).Info("failed to check user and password")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	permissionsString, err := getPermissionsString(PermissionsQ(r), user.Status)
	if err != nil {
		Log(r).WithError(err).Info("failed to create permissions string")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	access, refresh, claims, err := generateTokens(data.GenerateTokens{
		User:              *user,
		AccessLife:        helpers.ParseDurationStringToUnix(JwtParams(r).AccessLife),
		Secret:            JwtParams(r).Secret,
		PermissionsString: permissionsString,
	})
	if err != nil {
		Log(r).WithError(err).Info("failed to generate access and refresh tokens")
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
		Log(r).WithError(err).Error(err, "failed to create instance of refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, models.NewAuthTokenResponse(access, refresh))
}

func checkUserAndPassword(request requests.LoginRequest, usersQ data.Users) (*data.User, error) {
	user, err := usersQ.FilterByEmail(request.Data.Attributes.Email).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	if user == nil {
		return nil, errors.Errorf("no user with such email `%s`", request.Data.Attributes.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Data.Attributes.Password))
	if err != nil {
		return nil, errors.Wrap(err, "wrong password")
	}

	return user, nil
}

func getPermissionsString(permissionsQ data.Permissions, userStatus data.UserStatus) (string, error) {
	permissions, err := permissionsQ.WithModules().FilterByStatus(userStatus).Select()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user permissions")
	}

	permissionsString, err := helpers.CreatePermissionsString(permissions)
	if err != nil {
		return "", errors.Wrap(err, "failed to get create user permissions string")
	}

	return permissionsString, nil
}

func generateTokens(dataToGenerate data.GenerateTokens) (access, refresh string, claims *data.JwtClaims, err error) {
	access, err = helpers.GenerateAccessToken(dataToGenerate)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to create access token")
	}

	refresh, err, claims = helpers.GenerateRefreshToken(dataToGenerate)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to create refresh token")
	}

	return access, refresh, claims, err
}
