package yammer

// List the current user's networks.
func (c *Client) ListNetworks() (rv []Network, err error) {
	u := "https://www.yammer.com/api/v1/networks/current.json"
	err = decodeReq(c, u, &rv, map[string]string{})
	return
}
