/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ModulePermission struct {
	Key
	Attributes ModulePermissionAttributes `json:"attributes"`
}
type ModulePermissionResponse struct {
	Data     ModulePermission `json:"data"`
	Included Included         `json:"included"`
}

type ModulePermissionListResponse struct {
	Data     []ModulePermission `json:"data"`
	Included Included           `json:"included"`
	Links    *Links             `json:"links"`
}

// MustModulePermission - returns ModulePermission from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustModulePermission(key Key) *ModulePermission {
	var modulePermission ModulePermission
	if c.tryFindEntry(key, &modulePermission) {
		return &modulePermission
	}
	return nil
}
