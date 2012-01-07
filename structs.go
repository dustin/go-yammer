package yammer

import (
	oauth "github.com/alloy-d/goauth"
)

type User struct {
	SO                *string                  `json:"significant_other,omitempty"`
	Schools           []map[string]interface{} `json:"schools,omitempty"`
	Stats             map[string]int
	URLs              []string                 `json:"external_urls,omitempty"`
	URL               *string                  `json:"web_url,omitempty"`
	Avatar            *string                  `json:"mugshot_url,omitempty"`
	Location          *string                  `json:",omitempty"`
	Timezone          *string                  `json:",omitempty"`
	YURL              *string                  `json:"url,omitempty"`
	Domains           []string                 `json:"network_domains,omitempty"`
	Kids              *string                  `json:"kids_names,omitempty"`
	PreviousCompanies []map[string]interface{} `json:"previous_companies,omitempty"`
	FullName          *string                  `json:"full_name,omitempty"`
	Birthday          *string                  `json:"birth_date,omitempty"`
	Expertise         *string                  `json:",omitempty"`
	Summary           *string                  `json:",omitempty"`
	NetworkID         int                      `json:"network_id,omitempty"`
	Name              *string                  `json:",omitempty"`
	NetworkName       *string                  `json:"network_name,omitempty"`
	Interests         *string                  `json:",omitempty"`
	Contact           map[string]interface{}
	HireDate          *string `json:"hire_date,omitempty"`
	ID                int
	CanBroadcast      bool    `json:"can_broadcast,string"`
	Admin             bool    `json:"admin,string"`
	JobTitle          *string `json:"job_title,omitempty"`
}

type Client struct {
	oauth oauth.OAuth
}
