package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dustin/yammer.go"

	oauth "github.com/dustin/goauth"
)

var key, secret, filename string

func init() {
	flag.StringVar(&key, "key", "", "consumer key")
	flag.StringVar(&secret, "secret", "", "consumer secret")
	flag.StringVar(&filename, "authfile", "../.auth", "output file name")
}

func main() {
	flag.Parse()
	if err := yammer.VerifyKeyAndSecret(key, secret); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	o := &oauth.OAuth{
		// Service:         "yammer",
		SignatureMethod: oauth.HMAC_SHA1,
		RequestTokenURL: "https://www.yammer.com/oauth/request_token",
		AccessTokenURL:  "https://www.yammer.com/oauth/access_token",
		OwnerAuthURL:    "https://www.yammer.com/oauth/authorize",
		ConsumerKey:     key,
		ConsumerSecret:  secret,
		Callback:        "",
	}

	if err := o.GetRequestToken(); err != nil {
		log.Fatalf("Error getting request token: %v", err)
	}

	url, err := o.AuthorizationURL()
	if err != nil {
		log.Fatalf("Error getting auth URL: %v", err)
	}

	log.Printf("Go get a PIN from %v", url)

	in := bufio.NewReader(os.Stdin)

	verifier, _, err := in.ReadLine()
	if err != nil {
		log.Fatalf("Error reading line:  %v", err)
	}

	err = o.GetAccessToken(string(verifier))
	if err != nil {
		log.Fatalf("Error verifying:  %v", err)
	}

	log.Printf("Auth junk:\n%#v", o)

	if err = o.Save(filename); err != nil {
		log.Fatalf("Error saving credentials.")
	}

	log.Printf("Credentials saved.")
}
