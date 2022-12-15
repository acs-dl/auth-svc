package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"gitlab.com/distributed_lab/acs/auth/internal/service/handlers"
	"gitlab.com/distributed_lab/acs/auth/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/auth/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

type Body struct {
	Data resources.Refresh
}

type Permissions struct {
	Data string `json:"data"`
}

func Jwt() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body Body

			bodyCopy, _ := io.ReadAll(r.Body)
			if err := json.NewDecoder(io.NopCloser(bytes.NewBuffer(bodyCopy))).Decode(&body.Data); err != nil {
				handlers.Log(r).WithError(err).Error(" failed to unmarshal")
				ape.RenderErr(w, problems.BadRequest(err)...)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))

			_, err := helpers.ParseJwtToken(body.Data.Attributes.Token, handlers.JwtParams(r).Secret)
			if err != nil {
				handlers.Log(r).WithError(err).Error("failed to decode jwt token")
				ape.RenderErr(w, problems.BadRequest(err)...)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
