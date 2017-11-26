package chatworkhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Hook is an inbound ChatWork webhook
type Hook struct {
	Signature  string
	RawPayload []byte
	Payload    WebhookPayload
}

// SignedBy checks that the provided secret matches the hook Signature
//
// Implements validation described in documentation:
// http://developer.chatwork.com/ja/webhook.html#requestSign
// http://creators-note.chatwork.com/entry/2017/11/22/165516
func (h *Hook) SignedBy(secret []byte) error {
	key, err := base64.StdEncoding.DecodeString(string(secret))
	if err != nil {
		return err
	}

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(h.RawPayload))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if h.Signature != expectedSignature {
		return errors.New("invalid signature")
	}

	return nil
}

// New reads a Hook from an incoming HTTP Request.
// Consider using `Parse` method rather than `New` method.
func New(req *http.Request) (hook *Hook, err error) {
	hook = new(Hook)
	if !strings.EqualFold(req.Method, "POST") {
		return nil, errors.New("unknown method")
	}

	if hook.Signature = req.Header.Get("X-ChatWorkWebhookSignature"); len(hook.Signature) == 0 {
		return nil, errors.New("no signature")
	}

	hook.RawPayload, err = ioutil.ReadAll(req.Body)
	return
}

// Parse reads and verifies the hook in an inbound request.
func Parse(secret []byte, req *http.Request) (hook *Hook, err error) {
	hook, err = New(req)
	if err != nil {
		return
	}

	err = hook.SignedBy(secret)

	var payload WebhookPayload
	err = json.Unmarshal(hook.RawPayload, &payload)
	hook.Payload = payload

	return
}

// EventType represents
type EventType int

const (
	// MessageCreated represents "message_created" webhook event
	MessageCreated = iota
	// MessageUpdated represents "message_updated" webhook event
	MessageUpdated
	// MentionToMe represents "mention_to_me" webhook event
	MentionToMe
)

// String stringer interface
func (e EventType) String() string {
	switch e {
	case MessageCreated:
		return "message_created"
	case MessageUpdated:
		return "message_updated"
	case MentionToMe:
		return "mention_to_me"
	default:
		return "unknown"
	}
}

// MarshalJSON marshal EventType
func (e EventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON unmarshal EventType
func (e *EventType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("data should be a string, got %s", data)
	}

	var et EventType
	switch s {
	case "message_created":
		et = MessageCreated
	case "message_updated":
		et = MessageUpdated
	case "mention_to_me":
		et = MentionToMe
	default:
		return fmt.Errorf("invalid EventType %s", s)

	}
	*e = et
	return nil
}

// Time represents
type Time struct {
	time.Time
}

// MarshalJSON marshal time to epoch time
func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

// UnmarshalJSON unmarshal epoch time to time
func (t *Time) UnmarshalJSON(data []byte) error {
	var s int64
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("data should be number, got %s", data)
	}
	*t = Time{time.Unix(s, 0)}
	return nil
}

// WebhookEvent represents event body
type WebhookEvent struct {
	MessageID  *string `json:"message_id,omitempty"`
	RoomID     *int    `json:"room_id,omitempty"`
	Body       *string `json:"body,omitempty"`
	SendTime   *Time   `json:"send_time,omitempty"`
	UpdateTime *Time   `json:"update_time,omitempty"`

	// AccountID may be nil when EventType is `MentionToMe`
	AccountID *int `json:"account_id,omitempty"`

	// FromAccountID may be nil when EventType are `MessageCreated` and `MessageUpdated`
	FromAccountID *int `json:"from_account_id,omitempty"`

	// ToAccountID may be nil when EventType are `MessageCreated` and `MessageUpdated`
	ToAccountID *int `json:"to_account_id,omitempty"`
}

// WebhookPayload represents webhook JSON
type WebhookPayload struct {
	SettingID string       `json:"webhook_setting_id,omitempty"`
	Type      EventType    `json:"webhook_event_type,omitempty"`
	Time      Time         `json:"webhook_event_time,omitempty"`
	Event     WebhookEvent `json:"webhook_event,omitempty"`
}
