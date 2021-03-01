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
	members, err := c.fetchConversationMembers(channelID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return c.filterActiveMember(members)
}

// CollectByUserGroup collect usergroup members using slack api.
func (c *MemberCollector) CollectByUserGroup(userGroupID string, channelID string) ([]Member, error) {
	ugMembers, err := c.client.GetUserGroupMembers(userGroupID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	chMembers, err := c.fetchConversationMembers(channelID)
	chMemberMap := make(map[string]string, len(chMembers))
	for _, chMember := range chMembers {
		chMemberMap[chMember] = chMember
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var members []string
	for _, ugMember := range ugMembers {
		if m, ok := chMemberMap[ugMember]; ok {
			members = append(members, m)
		}
	}
	return c.filterActiveMember(members)
}

func (c *MemberCollector) fetchConversationMembers(channelID string) ([]string, error) {
	nextCursor := ""
	params := slack.GetUsersInConversationParameters{
		ChannelID: channelID,
		Cursor:    nextCursor,
	}
	var members []string
	fetchedMembers, nextCursor, err := c.client.GetUsersInConversation(&params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	members = append(members, fetchedMembers...)
	for len(nextCursor) > 0 {
		fetchedMembers, nextCursor, err = c.client.GetUsersInConversation(&params)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		members = append(members, fetchedMembers...)
	}
	return members, err
}

func (c *MemberCollector) filterActiveMember(members []string) ([]Member, error) {
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
