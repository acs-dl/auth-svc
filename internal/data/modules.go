package data

type Modules interface {
	New() Modules

	Create(module ModulePermission) (ModulePermission, error)
	Select() ([]ModulePermission, error)
	Get() (*ModulePermission, error)
	Delete(moduleName string) error

	FilterByModule(moduleName string) Modules
	FilterById(id int64) Modules

	ResetFilters() Modules
}

type ModulePermission struct {
	Id         int64  `db:"id" structs:"-"`
	ModuleName string `db:"module_name" structs:"module_name"`
	Permission string `db:"permission" structs:"permission"`
}
