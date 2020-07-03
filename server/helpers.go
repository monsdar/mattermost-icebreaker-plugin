package main

import (
	"math/rand"

	"github.com/mattermost/mattermost-server/v5/model"
)

// GetRandomUser returns a random user that is found in the given channel and that is not a bot
// This function is limited to 1000 users per channel
func (p *Plugin) GetRandomUser(channelID string) (*model.User, *model.AppError) {
	//get a random user that is not a bot
	users, _ := p.API.GetUsersInChannel(channelID, "username", 0, 1000)
	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})

	targetuser := new(model.User)
	hasUserBeenFound := false
	for _, user := range users {
		if user.IsBot {
			continue
		}
		targetuser = user
		hasUserBeenFound = true
		break
	}

	if !hasUserBeenFound {
		return nil, &model.AppError{
			Message: "There is no user I can ask a question for...",
		}
	}
	return targetuser, nil
}
