package main

import (
	"log"

	"github.com/slack-go/slack"
)

// MemberCollector decide member randomly and make slackAttachment.
type MemberCollector struct {
	client *slack.Client
}

// Collect channnel members using slack api.
func (c *MemberCollector) Collect(channelID string) ([]Member, error) {
	chInfo, err := c.client.GetChannelInfo(channelID)
	// err return when channel is private. Try GetGroupInfo.
	if err != nil {
		grInfo, err := c.client.GetGroupInfo(channelID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return c.mapActiveMember(grInfo.Members)
	}
	return c.mapActiveMember(chInfo.Members)
}

func (c *MemberCollector) mapActiveMember(members []string) ([]Member, error) {
	var activeMembers []Member
	for _, mem := range members {
		user, err := c.client.GetUserInfo(mem)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if !(user.IsBot || user.Deleted) {
			activeMembers = append(activeMembers, Member{ID: user.ID, Name: user.Name})
		}
	}

	return activeMembers, nil
}
