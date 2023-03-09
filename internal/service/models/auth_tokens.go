package models

import "gitlab.com/distributed_lab/acs/auth/resources"

func NewAuthTokenResponse(access, refresh string) resources.AuthTokenResponse {
	return resources.AuthTokenResponse{
		Data: newAuthToken(access, refresh),
	}
}

func newAuthToken(access, refresh string) resources.AuthToken {
	return resources.AuthToken{
		Attributes: resources.AuthTokenAttributes{
			Access:  access,
			Refresh: refresh,
		},
	}
}
