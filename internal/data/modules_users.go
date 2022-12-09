package data

type ModulesUsers interface {
	New() ModulesUsers

	Create(moduleUser ModuleUser) (ModuleUser, error)
	Select() ([]ModuleUser, error)
	Delete(permission ModuleUser) error

	FilterByUserId(userId int64) ModulesUsers
	FilterByModuleId(moduleId int64) ModulesUsers
}

type ModuleUser struct {
	ModuleId int64 `db:"module_id" structs:"module_id"`
	UserId   int64 `db:"user_id" structs:"user_id"`
}
