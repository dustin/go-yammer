package yammer

import (
	"encoding/json"
	"errors"

	oauth "github.com/alloy-d/goauth"
)

func VerifyKeyAndSecret(key, secret string) error {
	if key == "" {
		return errors.New("Consumer key is required.")
	}
	if secret == "" {
		return errors.New("Consumer secret is required.")
	}
	return nil
}

// Get a client from the given auth file, consumer key, and secret
func New(authFile, key, secret string) (Client, error) {
	client := Client{
		oauth: oauth.OAuth{
			SignatureMethod: oauth.HMAC_SHA1,
			ConsumerKey:     key,
			ConsumerSecret:  secret,
		},
	}
	if err := client.oauth.Load(authFile); err != nil {
		return client, err
	}
	return client, nil
}

func decodeReq(c *Client, u string, rv interface{}) error {
	params := make(map[string]string)
	res, err := c.oauth.Get(u, params)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	defer res.Body.Close()

	d := json.NewDecoder(res.Body)
	if err = d.Decode(&rv); err != nil {
		return err
	}

	return nil
}

// Get the full list of users.
func (c *Client) ListUsers() ([]User, error) {
	u := "https://www.yammer.com/api/v1/users.json"

	rv := make([]User, 0)
	if err := decodeReq(c, u, &rv); err != nil {
		return rv, err
	}

	return rv, nil
}
