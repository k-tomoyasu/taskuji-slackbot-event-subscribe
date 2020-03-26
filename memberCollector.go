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

// CollectByUserGroup collect usergroup members using slack api.
func (c *MemberCollector) CollectByUserGroup(userGroupID string, channelID string) ([]Member, error) {
	ugMembers, err := c.client.GetUserGroupMembers(userGroupID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	chInfo, err := c.client.GetChannelInfo(channelID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var members []string
	for _, ugMember := range ugMembers {
		for _, chMember := range chInfo.Members {
			if ugMember == chMember {
				members = append(members, ugMember)
			}
		}
	}
	return c.mapActiveMember(members)
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
