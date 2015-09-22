# Github Webhook Server

package `gws` provides a `http.Server` implementation that handles webhook events from GitHub.

# Usage

```go
    s := gws.New("")

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

    log.Fatal(http.ListenAndServe(":8080", s))
```