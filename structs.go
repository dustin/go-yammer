// A yammer client.
package yammer

import (
	oauth "github.com/dustin/goauth"
)

// User as returned from ListUsers, etc...
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
	FullName          string                   `json:"full_name"`
	Birthday          *string                  `json:"birth_date,omitempty"`
	Expertise         *string                  `json:",omitempty"`
	Summary           *string                  `json:",omitempty"`
	NetworkID         int                      `json:"network_id,omitempty"`
	Name              string
	NetworkName       *string `json:"network_name,omitempty"`
	Interests         *string `json:",omitempty"`
	Contact           map[string]interface{}
	HireDate          *string `json:"hire_date,omitempty"`
	ID                int
	CanBroadcast      bool    `json:"can_broadcast,string"`
	Admin             bool    `json:"admin,string"`
	JobTitle          *string `json:"job_title,omitempty"`
}

// Group as returned from ListGroups, etc...
type Group struct {
	Privacy     string
	URL         string `json:"web_url"`
	Stats       map[string]interface{}
	Avatar      *string `json:"mugshot_url,omitempty"`
	YURL        *string `json:"url,omitempty"`
	Description *string
	FullName    string `json:"full_name"`
	Name        string
	ID          int
	CreatedAt   *string `json:"created_at"`
}

// Message as returned from a message list API
type Message struct {
	MessageType *string `json:"message_type"`
	LikedBy     struct {
		Count int
		Names []map[string]string
	}                 `json:"liked_by"`
	CreatedAt *string `json:"created_at"`
	ThreadId  int     `json:"thread_id"`
	APIURL    *string `json:"url"`
	NetworkID int     `json:"network_id"`
	Body      struct {
		Plain  string
		Parsed string
	}
	Privacy       *string
	ClientURL     *string `json:"client_url"`
	SenderType    *string `json:"sender_type"`
	ID            int
	SystemMessage bool `json:"system_message"`
	Attachments   []interface{}
	SenderId      int     `json:"sender_id"`
	RepliedTo     *int    `json:"replied_to_id"`
	ClientType    *string `json:"client_type"`
	DirectMessage bool    `json:"direct_message"`
	YURL          *string `json:"web_url"`
}

// Message request object for PostMessage
type MessageRequest struct {
	Body      string // Message body (required)
	GroupId   int    // The group to post to (optional)
	ReplyTo   int    // Message in reply to (optional)
	DirectTo  int    // ID of the user to whom this message is directed (optional)
	Broadcast bool   // True if an administrative broadcast
}

// The client.
type Client struct {
	oauth oauth.OAuth
}
