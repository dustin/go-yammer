package yammer

import (
	"strconv"
)

// Get the full list of users.
func (c *Client) ListUsers() ([]User, error) {
	u := "https://www.yammer.com/api/v1/users.json"

	rv := make([]User, 0)
	stillgoing := true

	for i := 1; stillgoing; i++ {
		m := make([]User, 0, 50)
		if err := decodeReq(c, u, &m,
			map[string]string{"page": strconv.FormatInt(int64(i), 10)}); err != nil {
			return rv, err
		}
		rv = append(rv, m...)
		stillgoing = len(m) == 50
	}

	return rv, nil
}
