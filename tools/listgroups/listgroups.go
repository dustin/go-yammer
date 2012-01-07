package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dustin/yammer.go"
)

var key, secret, filename string

func init() {
	flag.StringVar(&key, "key", "", "consumer key")
	flag.StringVar(&secret, "secret", "", "consumer secret")
	flag.StringVar(&filename, "authfile", "../.auth", "auth file path")
}

func main() {

	flag.Parse()
	if err := yammer.VerifyKeyAndSecret(key, secret); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	client, err := yammer.New(filename, key, secret)
	if err != nil {
		log.Fatalf("Error making client:  %v", err)
	}

	groups, err := client.ListGroups()
	if err != nil {
		log.Fatalf("Error listing groups:  %v", err)
	}

	b, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling groups:  %v", err)
	}
	os.Stdout.Write(b)
}
