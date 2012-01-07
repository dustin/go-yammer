package yammer

import (
	"errors"
	"strconv"
)

// Post a message.
func (c *Client) PostMessage(req MessageRequest) error {
	u := "https://www.yammer.com/api/v1/messages.json"

	params := make(map[string]string)
	params["body"] = req.Body
	if req.GroupId != 0 {
		params["group_id"] = strconv.FormatInt(int64(req.GroupId), 10)
	}
	if req.ReplyTo != 0 {
		params["reply_to_id"] = strconv.FormatInt(int64(req.ReplyTo), 10)
	}
	if req.DirectTo != 0 {
		params["direct_to_id"] = strconv.FormatInt(int64(req.DirectTo), 10)
	}
	if req.Broadcast {
		params["broadcast"] = "true"
	}
	res, err := c.oauth.Post(u, params)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		return errors.New(res.Status)
	}

	res.Body.Close()

	return nil
}
