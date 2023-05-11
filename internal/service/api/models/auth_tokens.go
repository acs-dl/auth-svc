package models

import "gitlab.com/distributed_lab/acs/auth/resources"

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
