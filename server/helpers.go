package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mroth/weightedrand"
)

// GetRandomUser returns a random user that is found in the given channel and that is not a bot
// This function is limited to 1000 users per channel
func (p *Plugin) GetRandomUser(channelID string, userIDToIgnore string) (*model.User, *model.AppError) {
	//get a random user that is not a bot
	users, _ := p.API.GetUsersInChannel(channelID, "username", 0, 1000)
	weightedUsers := []weightedrand.Choice{} //list of users, sorted by weight

	//read the users data for weightedrandom
	data := p.ReadFromStorage()

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

		//check if the user has already been asked lately. Add him with a weight according to how recent the asking has been
		//by iterating in reverse we make sure that users that appear multiple times in the list will not mess up the weights
		isNewUser := true
		if data.LastUsers != nil {
			for index := len(data.LastUsers) - 1; index >= 0; index-- {
				currentUserID := data.LastUsers[index]
				if currentUserID == user.Id {
					userWeight := uint(math.Abs(float64(index - len(data.LastUsers))))
					weightedUsers = append(weightedUsers, weightedrand.Choice{Weight: userWeight, Item: user})
					isNewUser = false
					break
				}
			}
			if !isNewUser { //if the user has been found within our data.LastUsers we can continue the loop
				continue
			}
		}

		//Finally... this is a brand-new user that has never asked a question. Add him with a very high weight, so he'll be chosen with a high possibility
		weightedUsers = append(weightedUsers, weightedrand.Choice{Weight: 1000, Item: user})
	}

	if len(weightedUsers) > 0 {
		chooser := weightedrand.NewChooser(weightedUsers...)
		user, ok := chooser.Pick().(*model.User)
		if ok {
			return user, nil
		}
	}

	return nil, &model.AppError{
		Message: "There is no user I can ask a question for...",
	}
}

// GetRandomQuestion returns a random question that hasn't been asked recently
func (p *Plugin) GetRandomQuestion() (*Question, *model.AppError) {
	weightedQuestions := []weightedrand.Choice{} //list of questions, sorted by weight

	//read the question data for weightedrandom
	data := p.ReadFromStorage()

	for _, question := range data.Questions {
		//check if the question has already been asked lately. Add it with a weight according to how recent the question has been asked
		//by iterating in reverse we make sure that questions that appear multiple times in the list will not mess up the weights
		isNewQuestion := true
		if data.LastQuestions != nil {
			for index := len(data.LastQuestions) - 1; index >= 0; index-- {
				currentQuestion := data.LastQuestions[index]
				if currentQuestion.Question == question.Question {
					questionWeight := uint(math.Abs(float64(index - len(data.LastQuestions))))
					weightedQuestions = append(weightedQuestions, weightedrand.Choice{Weight: questionWeight, Item: currentQuestion})
					isNewQuestion = false
					break
				}
			}
			if !isNewQuestion { //if the question has been found within our data.LastQuestions we can continue the loop
				continue
			}
		}

		//Finally... this is a brand-new question that has never been asked. Add it with a very high weight, so it'll be chosen with a high possibility
		weightedQuestions = append(weightedQuestions, weightedrand.Choice{Weight: 1000, Item: question})
	}

	if len(weightedQuestions) > 0 {
		chooser := weightedrand.NewChooser(weightedQuestions...)
		question, ok := chooser.Pick().(Question)
		if ok {
			return &question, nil
		}
	}

	return nil, &model.AppError{
		Message: "There is no question to ask...",
	}
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

	for _, field := range commandFields {
		index, err := strconv.Atoi(field)
		if err != nil {
			//the word we got is not a valid index, but perhaps the next fits...
			continue
		}
		if (len(givenArray) <= index) || (index < 0) {
			return -1, &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf("Error: Your given index of %d is not valid", index),
			}
		}
		return index, nil
	}

	return -1, &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Error: Please enter a valid index",
	}
}
