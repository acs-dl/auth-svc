package models

import (
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/resources"
)

func NewPermissionUserModel(PermissionUser data.PermissionUser) resources.PermissionUser {
	result := resources.PermissionUser{
		Key: resources.NewKeyInt64(PermissionUser.UserId, resources.MODULE_PERMISSION),
		Attributes: resources.PermissionUserAttributes{
			UserId:       PermissionUser.UserId,
			PermissionId: PermissionUser.PermissionId,
		},
	}

	return result
}

func NewPermissionUsersList(PermissionUsers []data.PermissionUser) []resources.PermissionUser {
	result := make([]resources.PermissionUser, len(PermissionUsers))
	for i, elem := range PermissionUsers {
		result[i] = NewPermissionUserModel(elem)
	}
	return result
}
