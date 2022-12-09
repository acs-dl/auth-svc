package models

import (
	"gitlab.com/distributed_lab/Auth/internal/data"
	"gitlab.com/distributed_lab/Auth/resources"
)

func NewModuleUserModel(moduleUser data.ModuleUser) resources.ModuleUser {
	result := resources.ModuleUser{
		Key: resources.NewKeyInt64(-1, resources.MODULE_PERMISSION),
		Attributes: resources.ModuleUserAttributes{
			UserId:   moduleUser.UserId,
			ModuleId: moduleUser.ModuleId,
		},
	}

	return result
}

func NewModuleUsersList(moduleUsers []data.ModuleUser) []resources.ModuleUser {
	result := make([]resources.ModuleUser, len(moduleUsers))
	for i, elem := range moduleUsers {
		result[i] = NewModuleUserModel(elem)
	}
	return result
}
