# chatworkhook

chatworkhook supports [ChatWork Webhook](http://developer.chatwork.com/ja/webhook.html).

## Installation

```sh
$ go get github.com/griffin-stewie/chatworkhook
```

## Usage

```go
webhookToken := []byte("<WEB HOOK TOKEN>")
hook, err := chatworkhook.Parse(webhookToken, request)
if err != nil {
    panic(err.Error())
}
body := *hook.Payload.Event.Body
fmt.Printf("Body is %s", body)
```
