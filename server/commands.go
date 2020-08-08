package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	commandIcebreaker               = "icebreaker"
	commandIcebreakerAdd            = commandIcebreaker + " add"
	commandIcebreakerList           = commandIcebreaker + " list"
	commandIcebreakerRemove         = commandIcebreaker + " admin remove"
	commandIcebreakerClearAll       = commandIcebreaker + " admin clearall"
	commandIcebreakerResetToDefault = commandIcebreaker + " admin reset questions"
)

func (p *Plugin) registerCommands() error {
	commands := [...]model.Command{
		model.Command{
			Trigger:          commandIcebreaker,
			AutoComplete:     true,
			AutoCompleteDesc: "Ask an icebreaker",
		},
		model.Command{
			Trigger:          commandIcebreakerAdd,
			AutoComplete:     true,
			AutoCompleteHint: "<question>",
			AutoCompleteDesc: "Propose as new icebreaker question",
		},
		model.Command{
			Trigger:          commandIcebreakerList,
			AutoComplete:     true,
			AutoCompleteDesc: "Show a list questions",
		},
		model.Command{
			Trigger:          commandIcebreakerRemove,
			AutoComplete:     true,
			AutoCompleteHint: "<id>",
			AutoCompleteDesc: "Remove a question. Admin only",
		},
		model.Command{
			Trigger:          commandIcebreakerClearAll,
			AutoComplete:     true,
			AutoCompleteDesc: "Remove ALL questions. Admin only",
		},
		model.Command{
			Trigger:          commandIcebreakerResetToDefault,
			AutoComplete:     true,
			AutoCompleteDesc: "Resets the questions to the default ones from this plugin. Admin only",
		},
	}

	for _, command := range commands {
		if err := p.API.RegisterCommand(&command); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("Failed to register %s command", commandIcebreakerResetToDefault))
		}
	}

	return nil
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand
// API.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	adminCommands := map[string]func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError){
		commandIcebreakerRemove: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerRemove(args), nil
		},
		commandIcebreakerClearAll: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerClearAll(args), nil
		},
		commandIcebreakerResetToDefault: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerResetToDefault(args), nil
		},
	}

	userCommands := map[string]func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError){
		commandIcebreakerAdd: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerAdd(args), nil
		},
		commandIcebreakerList: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerList(args), nil
		},
	}

	//this needs to be last, as prefix `/icebreaker` is also part of the above commands
	triggerCommand := func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
		return p.executeCommandIcebreaker(args), nil
	}

	trigger := strings.TrimPrefix(args.Command, "/")
	trigger = strings.TrimSuffix(trigger, " ")

	//first check for admin commands, make sure the user has the right permission
	for key, value := range adminCommands {
		if strings.HasPrefix(trigger, key) {
			sourceUser, _ := p.API.GetUser(args.UserId)
			if response := requireAdminUser(sourceUser); response != nil {
				return response, nil
			}
			return value(args)
		}
	}

	//then go for the user commands
	for key, value := range userCommands {
		if strings.HasPrefix(trigger, key) {
			return value(args)
		}
	}

	//last but not least check for the triggerCommand (it needs to be asked without any text behind it)
	if trigger == commandIcebreaker {
		return triggerCommand(args)
	}

	//return an error message when the command has not been detected at all
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Unknown command: " + args.Command),
	}, nil
}

func (p *Plugin) executeCommandIcebreakerResetToDefault(args *model.CommandArgs) *model.CommandResponse {
	p.FillDefaultQuestions()

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("All questions have been reset to the default ones. Beware the pitchforks!"),
	}
}

func (p *Plugin) executeCommandIcebreakerClearAll(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage()
	lenBefore := len(data.Questions)
	data.Questions = []Question{}
	p.WriteToStorage(&data)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("All %d proposed questions have been removed. Beware the pitchforks!", lenBefore),
	}
}

func (p *Plugin) executeCommandIcebreakerClearAllApproved(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage()
	lenBefore := len(data.Questions)
	data.Questions = []Question{}
	p.WriteToStorage(&data)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("All %d questions have been removed. Beware the pitchforks!", lenBefore),
	}
}

func (p *Plugin) executeCommandIcebreakerRemove(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage()

	indeces, errResponse := getIndeces(args.Command, data.Questions)
	if errResponse != nil {
		return errResponse
	}

	for _, index := range indeces {
		//from https://stackoverflow.com/a/37335777/199513
		data.Questions = append(data.Questions[:index], data.Questions[index+1:]...)
	}
	p.WriteToStorage(&data)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Questions removed",
	}
}

func (p *Plugin) executeCommandIcebreakerList(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage()

	if len(data.Questions) == 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "There are no questions...",
		}
	}

	message := "Questions:\n"
	for index, question := range data.Questions {
		creator := question.Creator
		user, err := p.API.GetUser(creator)
		if err == nil {
			creator = user.GetDisplayName("")
		}
		message = message + fmt.Sprintf("%d.\t%s:\t%s\n", index, creator, question.Question)
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         message,
	}
}

func (p *Plugin) executeCommandIcebreaker(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage()

	//check if there are any questions yet
	if len(data.Questions) == 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: There are no questions that I can ask. Be the first one to propose a question by using `/icebreaker add <question>`",
		}
	}

	//get a random user that is not a bot
	user, err := p.GetRandomUser(args.ChannelId, args.UserId)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: Cannot get a user to ask a question for. Note: This plugin will not ask questions to offline or DND users.",
		}
	}

	//build the question and ask it
	question, err := p.GetRandomQuestion()
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: There are no questions that I can ask. Be the first one to propose a question by using `/icebreaker add <question>`",
		}
	}

	message := fmt.Sprintf("Hey @%s! %s", user.GetDisplayName(""), question.Question)
	post := &model.Post{
		ChannelId: args.ChannelId,
		RootId:    args.RootId,
		UserId:    p.botID,
		Message:   message,
	}

	//store the user and question so we avoid asking the same users and same questions over and over
	data.LastUsers = append(data.LastUsers, user.Id)
	if len(data.LastUsers) > LenHistory {
		index := 0 //remove the oldest element
		data.LastUsers = append(data.LastUsers[:index], data.LastUsers[index+1:]...)
	}
	data.LastQuestions = append(data.LastQuestions, *question)
	if len(data.LastQuestions) > LenHistory {
		index := 0 //remove the oldest element
		data.LastQuestions = append(data.LastQuestions[:index], data.LastQuestions[index+1:]...)
	}
	p.WriteToStorage(&data)

	if _, err = p.API.CreatePost(post); err != nil {
		const errorMessage = "Error: Failed to create post"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         errorMessage,
		}
	}

	return &model.CommandResponse{}
}

func (p *Plugin) executeCommandIcebreakerAdd(args *model.CommandArgs) *model.CommandResponse {
	//check the user input and extract the question from it
	givenQuestion := strings.TrimPrefix(args.Command, fmt.Sprintf("/%s", commandIcebreakerAdd))
	givenQuestion = strings.TrimPrefix(givenQuestion, " ")
	if len(givenQuestion) <= 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: Please enter a question",
		}
	}

	newQuestion := Question{}
	creator, _ := p.API.GetUser(args.UserId)
	newQuestion.Creator = creator.Id
	newQuestion.Question = givenQuestion

	data := p.ReadFromStorage()

	//Check if the question is already created
	for _, question := range data.Questions {
		if question.Question == newQuestion.Question {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "Error: Your question has already been added",
			}
		}
	}
	data.Questions = append(data.Questions, newQuestion)

	p.WriteToStorage(&data)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Thanks %s! Added your question: '%s'. Total number of questions: %d", creator.GetDisplayName(""), newQuestion.Question, len(data.Questions)),
	}
}
