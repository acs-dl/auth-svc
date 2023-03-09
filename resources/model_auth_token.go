/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type AuthToken struct {
	Key
	Attributes AuthTokenAttributes `json:"attributes"`
}
type AuthTokenResponse struct {
	Data     AuthToken `json:"data"`
	Included Included  `json:"included"`
}

type AuthTokenListResponse struct {
	Data     []AuthToken `json:"data"`
	Included Included    `json:"included"`
	Links    *Links      `json:"links"`
}

// MustAuthToken - returns AuthToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustAuthToken(key Key) *AuthToken {
	var authToken AuthToken
	if c.tryFindEntry(key, &authToken) {
		return &authToken
	}
	return nil
}
