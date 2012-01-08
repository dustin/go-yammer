package yammer

import (
	"sort"
	"strconv"
	"strings"
)

type userSlice []User

func (u userSlice) Len() int {
	return len(u)
}

func (u userSlice) Less(i, j int) bool {
	return strings.ToLower(u[i].FullName) < strings.ToLower(u[j].FullName)
}

func (u userSlice) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

// Get the full list of users.
func (c *Client) ListUsers() ([]User, error) {
	u := "https://www.yammer.com/api/v1/users.json"

	rv := make(userSlice, 0)
	stillgoing := true

	for i := 1; stillgoing; i++ {
		m := make(userSlice, 0, 50)
		if err := decodeReq(c, u, &m,
			map[string]string{"page": strconv.FormatInt(int64(i), 10)}); err != nil {
			return rv, err
		}
		rv = append(rv, m...)
		stillgoing = len(m) == 50
	}

	sort.Sort(rv)

	return []User(rv), nil
}
