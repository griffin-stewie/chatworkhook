package chatworkhook_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/griffin-stewie/chatworkhook"
)

/// sample data from http://creators-note.chatwork.com/entry/2017/11/22/165516
const testToken = "A9ne+ygvdV0IZBaPFV2zC1e5Bk+IsI14BPwieRoBQNU="

func TestValidSignature(t *testing.T) {

	body := `{"webhook_setting_id":"246","webhook_event_type":"message_created","webhook_event_time":1511238729,"webhook_event":{"message_id":"984676321621704704","room_id":36818150,"account_id":1484814,"body":"test","send_time":1511238729,"update_time":0}}`
	signature := `G7Gtrh5Ee6d8erOVXhWPtUrkNJqqIT5vwLU50KhyLQk=`

	r, _ := http.NewRequest("POST", "/path", strings.NewReader(body))
	r.Header.Add("X-ChatWorkWebhookSignature", signature)

	if _, err := chatworkhook.Parse([]byte(testToken), r); err != nil {
		t.Error(fmt.Sprintf("Unexpected error '%s'", err))
	}
}

func TestEventParse(t *testing.T) {

	body := `{"webhook_setting_id":"246","webhook_event_type":"message_created","webhook_event_time":1511238729,"webhook_event":{"message_id":"984676321621704704","room_id":36818150,"account_id":1484814,"body":"test","send_time":1511238729,"update_time":0}}`
	signature := `G7Gtrh5Ee6d8erOVXhWPtUrkNJqqIT5vwLU50KhyLQk=`

	r, _ := http.NewRequest("POST", "/path", strings.NewReader(body))
	r.Header.Add("X-ChatWorkWebhookSignature", signature)

	h, err := chatworkhook.Parse([]byte(testToken), r)
	if err != nil {
		t.Errorf("Unexpected error '%s'", err)
	}

	payload := h.Payload

	if payload.Type != chatworkhook.MessageCreated {
		t.Errorf("EventType should be 'MessageCreated'")
	}

	if payload.SettingID != "246" {
		t.Errorf("Invalid SettingID, expected '246'")
	}

	if payload.Time.Time != time.Unix(1511238729, 0) {
		t.Errorf("Invalid EventTime, expected 1511238729")
	}

	event := payload.Event

	if event.SendTime.Time != time.Unix(1511238729, 0) {
		t.Errorf("Invalid EventTime, expected 1511238729")
	}

	if event.UpdateTime.Time != time.Unix(0, 0) {
		t.Errorf("Invalid EventTime, expected 0")
	}

	if *event.Body != "test" {
		t.Errorf("Invalid Body, expected 'test'")
	}

	if *event.RoomID != 36818150 {
		t.Errorf("Invalid RoomID, expected 1484814")
	}

	if *event.AccountID != 1484814 {
		t.Errorf("Invalid AccountID, expected 1484814")
	}

	if *event.MessageID != "984676321621704704" {
		t.Errorf("Invalid AccountID, expected '984676321621704704")
	}
}
