/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ModuleUser struct {
	Key
	Attributes ModuleUserAttributes `json:"attributes"`
}
type ModuleUserResponse struct {
	Data     ModuleUser `json:"data"`
	Included Included   `json:"included"`
}

type ModuleUserListResponse struct {
	Data     []ModuleUser `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustModuleUser - returns ModuleUser from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustModuleUser(key Key) *ModuleUser {
	var moduleUser ModuleUser
	if c.tryFindEntry(key, &moduleUser) {
		return &moduleUser
	}
	return nil
}
