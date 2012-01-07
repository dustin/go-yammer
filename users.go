package yammer

// Get the full list of users.
func (c *Client) ListUsers() ([]User, error) {
	u := "https://www.yammer.com/api/v1/users.json"

	rv := make([]User, 0)
	if err := decodeReq(c, u, &rv, map[string]string{}); err != nil {
		return rv, err
	}

	return rv, nil
}
