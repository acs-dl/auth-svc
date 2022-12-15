package helpers

import (
	"fmt"

	"gitlab.com/distributed_lab/acs/auth/internal/data"
)

func CreatePermissionsString(permissions []data.ModuleUser, ModulesQ data.Modules) (string, error) {
	var resultPermission string

	for _, permission := range permissions {
		module, err := ModulesQ.FilterById(permission.ModuleId).Get()
		if err != nil {
			return "", err
		}
		resultPermission += fmt.Sprintf("%s.%s/", module.ModuleName, module.Permission)
		ModulesQ.ResetFilters()
	}
	return resultPermission, nil
}
