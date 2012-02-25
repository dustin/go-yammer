package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log/syslog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/yammer.go"
)

var key, secret, filename, addr string
var debug, transmit, logSyslog bool

var client yammer.Client

type inputxml struct {
	ID          int    `xml:"id"`
	Version     int    `xml:"version"`
	EventType   string `xml:"event_type"`
	Occurred    string `xml:"occurred_at"`
	Author      string `xml:"author"`
	ProjectId   int    `xml:"project_id"`
	Description string `xml:"description"`
	Stories     []struct {
		XMLName   xml.Name `xml:"story"`
		ID        int      `xml:"id"`
		URL       string   `xml:"url"`
		Name      string   `xml:"name"`
		StoryType string   `xml:"story_type"`
		State     string   `xml:"current_state"`
		Owner     string   `xml:"owned_by"`
		Requestor string   `xml:"requested_by"`
	} `xml:"stories>story"`
}

var log io.Writer

func init() {

	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.BoolVar(&transmit, "xmit", true, "enable transmitting (disable for debug)")
	flag.BoolVar(&logSyslog, "syslog", false, "log to syslog")
	flag.StringVar(&key, "key", "", "consumer key")
	flag.StringVar(&secret, "secret", "", "consumer secret")
	flag.StringVar(&filename, "authfile", "../.auth", "auth file path")

	flag.StringVar(&addr, "addr", ":8888", "web bind address")
}

func logF(format string, v ...interface{}) {
	fmt.Fprintf(log, format, v...)
	if !logSyslog {
		log.Write([]byte("\n"))
	}
}

func parseInt(params url.Values, key string) (rv int) {
	s := params.Get(key)
	if s != "" {
		i, err := strconv.ParseInt(params.Get(key), 10, 0)
		if err != nil {
			logF("Error parsing param %s: %v", key, err)
			return
		}
		rv = int(i)
	}
	return
}

func debugLog(input inputxml, output yammer.MessageRequest) {
	inputj, e := json.MarshalIndent(input, "", "  ")
	if e != nil {
		logF("Error generating input JSON:  %v", e)
		return
	}
	outputj, e := json.MarshalIndent(output, "", "  ")
	if e != nil {
		logF("Error generating output JSON:  %v", e)
		return
	}

	logF("Input:\n%s\nOutput:\n%s\n", inputj, outputj)

}

func splitWrite(w io.Writer) io.Writer {
	return io.MultiWriter(w, log)
}

func yammerPoster(w http.ResponseWriter, req *http.Request) {
	input := inputxml{}
	d := xml.NewDecoder(req.Body)
	if err := d.Decode(&input); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(splitWrite(w), "Error parsing input:  %v\n", err)
		return
	}

	params := req.URL.Query()

	tags := []string{fmt.Sprintf("#pt_%s", input.EventType)}

	if tag := params.Get("tag"); tag != "" {
		tags = append(tags, fmt.Sprintf("#%s", tag))
		tags = append(tags, fmt.Sprintf("#%s_%s", tag, input.EventType))
	}

	// project ID, story ID
	msg := fmt.Sprintf("%s %s\nhttps://www.pivotaltracker.com/projects/%d?story_id=%d",
		input.Description, strings.Join(tags, " "),
		input.ProjectId, input.Stories[0].ID)

	yreq := yammer.MessageRequest{
		Body:     msg,
		GroupId:  parseInt(params, "group_id"),
		ReplyTo:  parseInt(params, "reply_to"),
		DirectTo: parseInt(params, "direct_to"),
	}

	if debug {
		go debugLog(input, yreq)
	} else {
		logF("Yammer msg: %s", yreq.Body)
	}

	if transmit {
		if err := client.PostMessage(yreq); err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(splitWrite(w), "Error posting message:  %v\n", err)
			return
		}
	} else {
		logF("[transmission disabled]")
	}

	w.WriteHeader(201)
}

func groupLister(w http.ResponseWriter, req *http.Request) {
	groups, err := client.ListGroups()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(splitWrite(w), "Problem listing groups:  %v", err)
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
		fmt.Fprintf(splitWrite(w),
			"Expected GET or POST, got %s\n", req.Method)
		return
	case "POST":
		yammerPoster(w, req)
	case "GET":
		groupLister(w, req)
	}
}

func main() {
	flag.Parse()

	if logSyslog {
		var err error
		log, err = syslog.New(syslog.LOG_INFO, "pivotal")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't initialize syslog:  %v", err)
			os.Exit(1)
		}
	} else {
		log = os.Stdout
	}

	var err error
	if err = yammer.VerifyKeyAndSecret(key, secret); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	client, err = yammer.NewFromFile(filename, key, secret)
	if err != nil {
		logF("Error making client:  %v", err)
		os.Exit(1)
	}

	http.HandleFunc("/", yammerHandler)
	logF("Listening on %s", addr)
	if !transmit {
		logF("[transmission disabled]")
	}
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		logF("ListenAndServe: ", err)
		os.Exit(1)
	}
}
