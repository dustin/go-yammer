package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/dustin/yammer.go"
)

var key, secret, filename, addr string

var client yammer.Client

type inputxml struct {
	ID          int
	Version     int
	Occurred    string `xml:"occurred_at>"`
	Author      string
	ProjectId   int `xml:"project_id>"`
	Description string
	Stories     []struct {
		XMLName   xml.Name `xml:"story"`
		ID        int
		URL       string
		Name      string
		StoryType string `xml:"story_type>"`
		State     string `xml:"current_state>"`
		Owner     string `xml:"owned_by>"`
		Requestor string `xml:"requested_by>"`
	} `xml:"stories>story"`
}

func init() {
	flag.StringVar(&key, "key", "", "consumer key")
	flag.StringVar(&secret, "secret", "", "consumer secret")
	flag.StringVar(&filename, "authfile", "../.auth", "auth file path")

	flag.StringVar(&addr, "addr", ":8888", "web bind address")
}

func parseInt(params url.Values, key string) int {
	i, err := strconv.ParseInt(params.Get(key), 10, 0)
	if err != nil {
		log.Printf("Error parsing param %s: %v", key, err)
		return 0
	}
	return int(i)
}

func yammerPoster(w http.ResponseWriter, req *http.Request) {
	input := inputxml{}
	if err := xml.Unmarshal(req.Body, &input); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Error parsing input:  %v", err)
		return
	}
	inputj, e := json.MarshalIndent(input, "", "  ")
	if e != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error generating input JSON:  %v", e)
		return
	}

	params := req.URL.Query()

	tag := params.Get("tag")
	if tag != "" {
		tag = fmt.Sprintf(" #%s", tag)
	}

	// project ID, story ID
	msg := fmt.Sprintf("%s%s\nhttps://www.pivotaltracker.com/projects/%d?story_id=%d",
		input.Description, tag, input.ProjectId, input.Stories[0].ID)

	yreq := yammer.MessageRequest{
		Body:     msg,
		GroupId:  parseInt(params, "group_id"),
		ReplyTo:  parseInt(params, "reply_to"),
		DirectTo: parseInt(params, "direct_to"),
	}

	outputj, e := json.MarshalIndent(yreq, "", "  ")
	if e != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error generating output JSON:  %v", e)
		return
	}

	log.Printf("Input:\n%s\nOutput:\n%s\n", inputj, outputj)

	if err := client.PostMessage(yreq); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error posting message:  %v", err)
		return
	}

	w.WriteHeader(201)
}

func groupLister(w http.ResponseWriter, req *http.Request) {
	groups, err := client.ListGroups()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Problem listing groups:  %v", err)
		return
	}

	header := []byte(`<html>
<head>
  <title>Group List</title>
  <style type="text/css">
    dt { font-weight: bold }
    dd { font-family: monospace }
  </style>
</head>
<body>
  <h1>List of Groups and their IDs</h1>
<dl>`)
	footer := []byte(`</dl></body></html>`)

	w.Header().Set("content-type", "text/html")
	w.Write(header)
	defer w.Write(footer)

	for _, g := range groups {
		msg := fmt.Sprintf("<dt>%s</dt><dd>%d</dd>\n", g.FullName, g.ID)
		w.Write([]byte(msg))
	}
}

func yammerHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Excpected GET or POST, got %s", req.Method)
		return
	case "POST":
		yammerPoster(w, req)
	case "GET":
		groupLister(w, req)
	}
}

func main() {
	flag.Parse()

	var err error
	if err = yammer.VerifyKeyAndSecret(key, secret); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	client, err = yammer.New(filename, key, secret)
	if err != nil {
		log.Fatalf("Error making client:  %v", err)
	}

	http.HandleFunc("/", yammerHandler)
	log.Printf("Listening on %s", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
