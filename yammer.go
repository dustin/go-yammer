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

func decodeReq(c *Client, u string, rv interface{},
	params map[string]string) error {

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
