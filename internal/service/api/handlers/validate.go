package handlers

import (
	"net/http"

	"github.com/acs-dl/auth-svc/internal/service/api/helpers"
	"github.com/acs-dl/auth-svc/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Validate(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewValidateRequest(r)
	if err != nil {
		Log(r).WithError(err).Errorf("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	refreshToken, err := RefreshTokensQ(r).FilterByTokens(request.Token).Get()
	if err != nil {
		Log(r).WithError(err).Error(err, "failed to get refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if refreshToken == nil {
		Log(r).Errorf("no token was found in db")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	_, err = helpers.CheckTokenValidity(refreshToken.Token, JwtParams(r).Secret)
	if err != nil {
		Log(r).WithError(err).Errorf("something wrong with refresh token")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ape.Render(w, http.StatusOK)
}
