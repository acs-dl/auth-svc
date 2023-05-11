package models

import "github.com/acs-dl/auth-svc/resources"

func NewAuthTokenResponse(access string) resources.AuthTokenResponse {
	return resources.AuthTokenResponse{
		Data: newAuthToken(access),
	}
}

func newAuthToken(access string) resources.AuthToken {
	return resources.AuthToken{
		Attributes: resources.AuthTokenAttributes{
			Access: access,
		},
	}
}
