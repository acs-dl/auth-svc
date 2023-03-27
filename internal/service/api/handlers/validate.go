package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/service/api/helpers"
	"gitlab.com/distributed_lab/acs/auth/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Validate(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewValidateRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	refreshToken, err := RefreshTokensQ(r).FilterByToken(request.Token).Get()
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

	_, err = helpers.CheckTokenValidity(refreshToken.Token, JwtParams(r).Secret)
	if err != nil {
		Log(r).WithError(err).Info("something wrong with refresh token")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ape.Render(w, http.StatusOK)
}
