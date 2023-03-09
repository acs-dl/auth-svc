package models

import (
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/resources"
)

func NewModulePermissionModel(module data.Module, permission data.Permission) resources.ModulePermission {
	result := resources.ModulePermission{
		Key: resources.NewKeyInt64(module.Id, resources.MODULE_PERMISSION),
		Attributes: resources.ModulePermissionAttributes{
			Module:     module.Name,
			Permission: permission.Name,
		},
	}

	return result
}

func NewModulePermissionsList(modulePermissions []data.ModulePermission) []resources.ModulePermission {
	result := make([]resources.ModulePermission, len(modulePermissions))
	for i, item := range modulePermissions {
		module := data.Module{
			Name: item.ModuleName,
			Id:   item.Id,
		}
		permission := data.Permission{
			Name:     item.PermissionName,
			ModuleId: item.ModuleId,
		}
		result[i] = NewModulePermissionModel(module, permission)
	}
	return result
}
