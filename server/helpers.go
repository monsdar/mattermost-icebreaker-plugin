package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

// GetRandomUser returns a random user that is found in the given channel and that is not a bot
// This function is limited to 1000 users per channel
func (p *Plugin) GetRandomUser(channelID string, userIDToIgnore string) (*model.User, *model.AppError) {
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
		if user.Id == userIDToIgnore {
			continue
		}
		status, err := p.API.GetUserStatus(user.Id)
		if (err != nil) || (status.Status == "offline") || (status.Status == "dnd") {
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

func requireAdminUser(sourceUser *model.User) *model.CommandResponse {
	if !sourceUser.IsSystemAdmin() { //TODO: Check for Channel owner instead of System Admin
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: You need to be admin in order to clear all proposed questions",
		}
	}
	return nil
}

func getIndex(command string, givenArray []Question) (int, *model.CommandResponse) {
	commandFields := strings.Fields(command)
	if len(commandFields) <= 2 {
		return 0, &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: Please enter a valid index",
		}
	}
	indexStr := commandFields[2]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return 0, &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error: Your given index of %s is not valid", indexStr),
		}
	}
	if len(givenArray) <= index {
		return 0, &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error: Your given index of %s is not valid", indexStr),
		}
	}

	return index, nil
}
