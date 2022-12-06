package handlers

import (
	"fmt"
	"github.com/mhrynenko/jwt_service/internal/data"
	"github.com/mhrynenko/jwt_service/internal/service/helpers"
	"github.com/mhrynenko/jwt_service/internal/service/requests"
	"github.com/mhrynenko/jwt_service/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Data.Attributes.Password)) != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	access, err := helpers.GenerateAccessToken(*user)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to create access token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	refresh, err, claims := helpers.GenerateRefreshToken(*user)
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to create refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	fmt.Println(claims["exp"])
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
