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

type listerFunc func(c *yammer.Client) (interface{}, error)

func listTool(key, secret, filename string, Lister listerFunc) {
	if err := yammer.VerifyKeyAndSecret(key, secret); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	client, err := yammer.New(filename, key, secret)
	if err != nil {
		log.Fatalf("Error making client:  %v", err)
	}

	stuff, err := Lister(&client)
	if err != nil {
		log.Fatalf("Error listing :  %v", err)
	}

	b, err := json.MarshalIndent(stuff, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling stuff:  %v", err)
	}
	os.Stdout.Write(b)
}

func main() {
	flag.Parse()

	funcs := map[string]listerFunc{
		"groups": func(c *yammer.Client) (interface{}, error) { return c.ListGroups() },
		"users":  func(c *yammer.Client) (interface{}, error) { return c.ListUsers() },
	}

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "List what?  I know of the following:\n")
		for k, _ := range funcs {
			fmt.Fprintf(os.Stderr, " - %s\n", k)
		}
		os.Exit(1)
	}

	f, ok := funcs[flag.Arg(0)]
	if !ok {
		log.Fatalf("I don't know how to list '%s'", flag.Arg(0))
	}

	listTool(key, secret, filename, f)
}
