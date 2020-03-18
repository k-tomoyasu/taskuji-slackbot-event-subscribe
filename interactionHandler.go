package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/slack-go/slack"
)

// interactionHandler handles interactive message response.
type interactionHandler struct {
	slackClient       *slack.Client
	verificationToken string
	lot               *Lot
	memberCollector   *MemberCollector
}

func (h interactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept message from slack with valid token
	if message.Token != h.verificationToken {
		log.Printf("[ERROR] Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	action := message.ActionCallback.AttachmentActions[0]
	h.reply(action, message, w)
}

func (h interactionHandler) reply(action *slack.AttachmentAction, message slack.AttachmentActionCallback, w http.ResponseWriter) {
	switch action.Name {
	case actionAccept:
		winnerReponsed := fmt.Sprintf("<@%s>", message.User.ID) == message.OriginalMessage.Text
		var value string
		if winnerReponsed {
			value = "Thank you:muscle:"
		} else {
			value = fmt.Sprintf("Oh,Thank you! <@%s>:muscle:", message.User.ID)
		}
		message.OriginalMessage.Attachments[0].Footer = "Good Luck!"
		responseMessage(w, message.OriginalMessage, "", value)
		return
	case actionRepeat:
		responseMessage(w, message.OriginalMessage, ":cry:", "")
		members, _ := h.memberCollector.Collect(message.Channel.ID)
		targetMembers := make([]Member, 0)
		// exclude member NG button pushed.
		for _, member := range members {
			if member.ID != message.User.ID {
				targetMembers = append(targetMembers, member)
			}
		}
		h.lot.DrawLots(message.Channel.ID, targetMembers)
		return
	default:
		log.Printf("[ERROR] ]Invalid action was submitted: %s", action.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
