package middleware

import (
	"net/http"
	"strings"

	"gitlab.com/distributed_lab/acs/auth/internal/service/handlers"
	"gitlab.com/distributed_lab/acs/auth/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Jwt(secret, module string, permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handlers.Log(r).Errorf("empty authorization header")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			splitAuthHeader := strings.Split(authHeader, " ")
			if len(splitAuthHeader) < 2 {
				//handlers.Log(r).Errorf("bad header structure")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}
			claims, err := helpers.ParseJwtToken(splitAuthHeader[1], secret)
			if err != nil {
				//handlers.Log(r).WithError(err).Error("failed to decode jwt token")
				ape.RenderErr(w, problems.BadRequest(err)...)
				return
			}

			splitModulePermission := strings.Split(claims["module.permission"].(string), "/")

			permissionMap := make(map[string]string)
			for _, modulePermission := range splitModulePermission {
				split := strings.Split(modulePermission, ".")
				if len(split) < 2 {
					continue
				}

				permissionMap[split[0]] = split[1]
			}

			for _, permission := range permissions {
				if permissionMap[module] == permission {
					next.ServeHTTP(w, r)
					return
				}
			}

			//handlers.Log(r).Errorf("allowed permission is higher than user's one")
			ape.RenderErr(w, problems.Forbidden())
		})
	}
}
