/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Login struct {
	Key
	Attributes LoginAttributes `json:"attributes"`
}
type LoginResponse struct {
	Data     Login    `json:"data"`
	Included Included `json:"included"`
}

type LoginListResponse struct {
	Data     []Login  `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustLogin - returns Login from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLogin(key Key) *Login {
	var login Login
	if c.tryFindEntry(key, &login) {
		return &login
	}
	return nil
}
