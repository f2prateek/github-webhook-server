package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/f2prateek/github-webhook-server"
)

var addr = flag.String("address", ":4001", "bind address")
var secret = flag.String("secret", "", "secret")

func main() {
	flag.Parse()

	s := gws.New(*secret)
	go func() {
		for {
			select {
			case event := <-s.PushEvents:
				fmt.Println("Received Push", event)
			case event := <-s.IssueEvents:
				fmt.Println("Received Issue", event)
			case event := <-s.IssueCommentEvents:
				fmt.Println("Received Issue Comment", event)
			case event := <-s.PullRequestEvents:
				fmt.Println("Received PR", event)
			case event := <-s.OtherEvents:
				fmt.Println("Received event", event)
			}
		}
	}()

	log.Printf("starting web server at %s with secret %s", *addr, *secret)
	log.Fatal(http.ListenAndServe(*addr, s))
}

func check(err error) {
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
