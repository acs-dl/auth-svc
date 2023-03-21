/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type PermissionUser struct {
	Key
	Attributes PermissionUserAttributes `json:"attributes"`
}
type PermissionUserResponse struct {
	Data     PermissionUser `json:"data"`
	Included Included       `json:"included"`
}

type PermissionUserListResponse struct {
	Data     []PermissionUser `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
}

// MustPermissionUser - returns PermissionUser from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPermissionUser(key Key) *PermissionUser {
	var permissionUser PermissionUser
	if c.tryFindEntry(key, &permissionUser) {
		return &permissionUser
	}
	return nil
}
