package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/slack-go/slack"
)

// interactionHandler handles interactive message response.
type interactionHandler struct {
	slackClient     *slack.Client
	signingSecret   string
	lot             *Lot
	memberCollector *MemberCollector
	messageTemplate MessageTemplate
}

func (h interactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Only accept message when verification success
	verifier, err := slack.NewSecretsVerifier(r.Header, h.signingSecret)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
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

	action := message.ActionCallback.AttachmentActions[0]
	h.reply(action, message, w)
}

func (h interactionHandler) reply(action *slack.AttachmentAction, message slack.AttachmentActionCallback, w http.ResponseWriter) {
	switch action.Name {
	case actionAccept:
		winnerResponsed := fmt.Sprintf("<@%s>", message.User.ID) == message.OriginalMessage.Text
		var value string
		if winnerResponsed {
			value = h.messageTemplate.WinnerResponded
		} else {
			value = fmt.Sprintf(h.messageTemplate.OtherResponded, message.User.ID)
		}
		message.OriginalMessage.Attachments[0].Footer = "Good Luck!"
		responseMessage(w, message.OriginalMessage, "", value)
		return
	case actionRepeat:
		responseMessage(w, message.OriginalMessage, ":cry:", "")
		userGroupID := action.Value
		var members []Member
		if len(userGroupID) > 0 {
			members, _ = h.memberCollector.CollectByUserGroup(userGroupID, message.Channel.ID)
		} else {
			members, _ = h.memberCollector.Collect(message.Channel.ID)
		}

		targetMembers := make([]Member, 0)
		// exclude member NG button pushed.
		for _, member := range members {
			if member.ID != message.User.ID {
				targetMembers = append(targetMembers, member)
			}
		}
		h.lot.DrawLots(message.Channel.ID, targetMembers, userGroupID)
		return
	default:
		log.Printf("[ERROR] ]Invalid action was submitted: %s", action.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
