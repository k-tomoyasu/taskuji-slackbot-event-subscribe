package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// eventHandler handles interactive message response.
type eventHandler struct {
	slackClient       *slack.Client
	verificationToken string
	lot               *Lot
	memberCollector   *MemberCollector
}

func (h eventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()

	event, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: h.verificationToken}))

	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if event.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}

	if event.Type == slackevents.CallbackEvent {
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			members, err := h.memberCollector.Collect(ev.Channel)
			if err != nil {
				log.Printf("[ERROR] Failed to collect member request: %s", err)
			}
			h.lot.DrawLots(ev.Channel, members, "")
		}
	}
}
