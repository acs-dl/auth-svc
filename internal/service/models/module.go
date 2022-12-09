package models

import (
	"gitlab.com/distributed_lab/Auth/internal/data"
	"gitlab.com/distributed_lab/Auth/resources"
)

func NewModulePermissionModel(module data.ModulePermission) resources.ModulePermission {
	result := resources.ModulePermission{
		Key: resources.NewKeyInt64(module.Id, resources.MODULE_PERMISSION),
		Attributes: resources.ModulePermissionAttributes{
			Module:     module.ModuleName,
			Permission: module.Permission,
		},
	}

	return result
}

func NewModulePermissionsList(modulePermissions []data.ModulePermission) []resources.ModulePermission {
	result := make([]resources.ModulePermission, len(modulePermissions))
	for i, elem := range modulePermissions {
		result[i] = NewModulePermissionModel(elem)
	}
	return result
}
