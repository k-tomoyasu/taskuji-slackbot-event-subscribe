package main

import (
	"encoding/json"
	"net/http"

	"github.com/slack-go/slack"
)

const (
	// action is used for slack attament action.
	actionAccept = "appcept"
	actionRepeat = "repeat"
)

const (
	// action type is used for slack attament action.
	repeatChannel = "repeatChannel"
	repeatGroup   = "repeatGroup"
)

// responseMessage response to the original slackbutton enabled message.
// It removes button and replace it with message which indicate how bot will work
func responseMessage(w http.ResponseWriter, original slack.Message, title, value string) {
	original.ReplaceOriginal = true
	original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}
