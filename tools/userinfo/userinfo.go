package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/dustin/yammer.go"
)

var key, secret, filename string
var shouldNotify bool

var wg sync.WaitGroup

var ignored map[string]bool = map[string]bool{
	"Github":          true,
	"Pivotal Tracker": true,
}

var notifyTemplateText = `
Hello {{.FullName}},

I noticed you don't have a {{.Field}} set up in your account.  Please go fill in a little about yourself for our company directory so that people can find you when we need you and all that.

You can find this and many more exciting fields here:

    https://www.yammer.com/account/profile_info

Love,

Your friendly yammer bot
(this is a recording)
`

var notifyTemplate = template.New("notify")

func init() {
	flag.StringVar(&key, "key", "", "consumer key")
	flag.StringVar(&secret, "secret", "", "consumer secret")
	flag.StringVar(&filename, "authfile", "../.auth", "auth file path")
	flag.BoolVar(&shouldNotify, "notify", false, "Should we notify the users?")

	var err error
	notifyTemplate, err = notifyTemplate.Parse(notifyTemplateText)
	if err != nil {
		log.Panicf("Error parsing template: %v", err)
	}
}

func phoneConvert(phoneX []interface{}) map[string]string {
	rv := make(map[string]string)
	for _, e := range phoneX {
		phone, ok := e.(map[string]interface{})
		if !ok {
			log.Fatalf("Error converting phone numbers from %#v", phoneX)
		}
		if len(phone) != 2 {
			log.Fatalf("Things are not as expected: %#v", phone)
		}
		rv[phone["type"].(string)] = phone["number"].(string)
	}
	return rv
}

func checkPhone(u yammer.User) (rv bool) {
	if phoneX, ok := u.Contact["phone_numbers"]; ok {
		phone, typed := phoneX.([]interface{})
		if !typed {
			log.Fatalf("Incorrect result in phone_numbers: %#v", phoneX)
		}
		if len(phone) > 0 {
			rv = true
		}
	}
	return
}

func checkIM(u yammer.User) (rv bool) {
	if imX, ok := u.Contact["im"]; ok {
		im, typed := imX.(map[string]interface{})
		if !typed {
			log.Fatalf("Incorrect result in im: %#v", imX)
		}
		if val, ok := im["username"]; ok && len(val.(string)) > 0 {
			rv = true
		}
	}
	return
}

func notifyMissingPhone(client yammer.Client, u yammer.User) {

	log.Printf("Notifying %v", u.FullName)

	b := bytes.NewBuffer(make([]byte, 0, 4096))
	data := map[string]string{
		"FullName": u.FullName,
		"Field":    "phone number",
	}
	err := notifyTemplate.Execute(b, data)
	if err != nil {
		log.Fatalf("Error executing template for %v: %v", u, err)
	}

	output := b.String()

	msg := yammer.MessageRequest{
		Body:      output,
		DirectTo:  u.ID,
		Broadcast: false,
	}

	for i := 0; i < 5; i++ {
		err = client.PostMessage(msg)

		if err == nil {
			return
		}
		log.Printf("Error posting message to %v: %v (retrying)",
			u.FullName, err)
		time.Sleep(time.Duration(i) * 5 * time.Second)
	}
}

func checkUser(client yammer.Client, u yammer.User) {
	defer wg.Done()
	if ignored[u.FullName] {
		return
	}
	if !checkPhone(u) {
		fmt.Printf("No phone numbers from %s\n", u.FullName)

		if shouldNotify {
			notifyMissingPhone(client, u)
		}
	}

	// if !checkIM(u) {
	// 	fmt.Printf("No IM from %s\n", u.FullName)
	// }
}

func dump(users []yammer.User) {
	f, err := os.Create("users.json")
	if err != nil {
		log.Fatalf("Error creating users file: %v", err)
	}
	defer f.Close()
	e := json.NewEncoder(f)
	err = e.Encode(users)
	if err != nil {
		log.Fatalf("Error serializing users: %v", err)
	}
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

	users, err := client.ListUsers()
	if err != nil {
		log.Fatalf("Error listing users: %v", err)
	}
	dump(users)

	for _, u := range users {
		wg.Add(1)
		go checkUser(client, u)
	}

	wg.Wait()
}
