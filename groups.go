package yammer

// Get the full list of users.
func (c *Client) ListGroups() ([]Group, error) {
	u := "https://www.yammer.com/api/v1/groups.json"

	rv := make([]Group, 0)
	if err := decodeReq(c, u, &rv, map[string]string{}); err != nil {
		return rv, err
	}

	return rv, nil
}
