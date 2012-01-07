package yammer

import (
	"sort"
	"strings"
)

type groupSlice []Group

func (g groupSlice) Len() int {
	return len(g)
}

func (g groupSlice) Less(i, j int) bool {
	return strings.ToLower(g[i].FullName) < strings.ToLower(g[j].FullName)
}

func (g groupSlice) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// Get the full list of users.
func (c *Client) ListGroups() ([]Group, error) {
	u := "https://www.yammer.com/api/v1/groups.json"

	rv := make(groupSlice, 0)
	if err := decodeReq(c, u, &rv, map[string]string{}); err != nil {
		return rv, err
	}

	sort.Sort(rv)

	return []Group(rv), nil
}
