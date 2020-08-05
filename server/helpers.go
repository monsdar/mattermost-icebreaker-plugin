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

	targetuser := new(model.User)
	hasUserBeenFound := false
	for len(users) > 0 { //as long as there are users left continue to search for a good candidate
		randomIndex := rand.Intn(len(users))
		currentUser := users[randomIndex]
		if currentUser.IsBot {
			users = append(users[:randomIndex], users[randomIndex+1:]...) //from https://stackoverflow.com/a/37335777/199513
			continue
		}
		if currentUser.Id == userIDToIgnore {
			users = append(users[:randomIndex], users[randomIndex+1:]...) //from https://stackoverflow.com/a/37335777/199513
			continue
		}
		status, err := p.API.GetUserStatus(currentUser.Id)
		if (err != nil) || (status.Status == "offline") || (status.Status == "dnd") {
			users = append(users[:randomIndex], users[randomIndex+1:]...) //from https://stackoverflow.com/a/37335777/199513
			continue
		}

		targetuser = currentUser
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

func getIndeces(command string, givenArray []Question) ([]int, *model.CommandResponse) {
	commandFields := strings.Fields(command)
	indeces := []int{}

	for _, field := range commandFields {
		index, err := strconv.Atoi(field)
		if err != nil {
			//do nothing... The word we got is not a valid index, but perhaps the next fits...
		}
		if len(givenArray) <= index {
			return []int{}, &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf("Error: Your given index of %d is not valid", index),
			}
		}
		indeces = append(indeces, index)
	}

	if len(indeces) == 0 {
		return []int{}, &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: Please enter a valid index",
		}
	}

	return indeces, nil
}
