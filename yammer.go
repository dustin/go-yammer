package yammer

import (
	"encoding/json"
	"errors"

	"github.com/dustin/goauth"
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

// Get a client using the given OAuth structure
func New(o oauth.OAuth) Client {
	return Client{oauth: o}
}

// Load OAuth data from a file.
func NewFromFile(authFile, key, secret string) (Client, error) {
	o := oauth.OAuth{
		SignatureMethod: oauth.HMAC_SHA1,
		ConsumerKey:     key,
		ConsumerSecret:  secret,
	}
	if err := o.Load(authFile); err != nil {
		return Client{}, err
	}
	return New(o), nil
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
