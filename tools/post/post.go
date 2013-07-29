package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dustin/go-yammer"
)

var key, secret, filename string
var groupId, replyTo, directTo int

func init() {
	flag.StringVar(&key, "key", "", "consumer key")
	flag.StringVar(&secret, "secret", "", "consumer secret")
	flag.StringVar(&filename, "authfile", "../.auth", "auth file path")

	flag.IntVar(&groupId, "groupId", 0, "group id (optional)")
	flag.IntVar(&replyTo, "replyTo", 0, "reply to (optional)")
	flag.IntVar(&directTo, "directTo", 0, "direct to (optional)")
}

func main() {
	flag.Parse()

	if err := yammer.VerifyKeyAndSecret(key, secret); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	client, err := yammer.NewFromFile(filename, key, secret)
	if err != nil {
		log.Fatalf("Error making client:  %v", err)
	}

	msg, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading message: %v", err)
	}

	req := yammer.MessageRequest{
		Body:     string(msg),
		GroupId:  groupId,
		ReplyTo:  replyTo,
		DirectTo: directTo,
	}

	err = client.PostMessage(req)
	if err != nil {
		log.Fatalf("Error posting message:  %v", err)
	}
}
