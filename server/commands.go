package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	commandIcebreaker                  = "icebreaker"
	commandIcebreakerAdd               = commandIcebreaker + " add"
	commandIcebreakerShowProposals     = commandIcebreaker + " show proposals"
	commandIcebreakerShowApproved      = commandIcebreaker + " show approved"
	commandIcebreakerApprove           = commandIcebreaker + " admin approve"
	commandIcebreakerReject            = commandIcebreaker + " admin reject"
	commandIcebreakerRemove            = commandIcebreaker + " admin remove"
	commandIcebreakerClearAllProposals = commandIcebreaker + " admin clearall proposals"
	commandIcebreakerClearAllApproved  = commandIcebreaker + " admin clearall approved"
	commandIcebreakerResetToDefault    = commandIcebreaker + " admin reset questions"
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
			Trigger:          commandIcebreakerApprove,
			AutoComplete:     true,
			AutoCompleteHint: "<id>",
			AutoCompleteDesc: "Approve a proposed IceBreaker question. Channel owners only",
		},
		model.Command{
			Trigger:          commandIcebreakerReject,
			AutoComplete:     true,
			AutoCompleteHint: "<id>",
			AutoCompleteDesc: "Reject a proposed IceBreaker question. Channel owners only",
		},
		model.Command{
			Trigger:          commandIcebreakerRemove,
			AutoComplete:     true,
			AutoCompleteHint: "<id>",
			AutoCompleteDesc: "Remove an already approved IceBreaker question. Channel owners only",
		},
		model.Command{
			Trigger:          commandIcebreakerClearAllProposals,
			AutoComplete:     true,
			AutoCompleteDesc: "Remove ALL proposed IceBreaker question. Channel owners only",
		},
		model.Command{
			Trigger:          commandIcebreakerClearAllApproved,
			AutoComplete:     true,
			AutoCompleteDesc: "Remove ALL approved IceBreaker question. Channel owners only",
		},
		model.Command{
			Trigger:          commandIcebreakerShowProposals,
			AutoComplete:     true,
			AutoCompleteDesc: "Show a list of proposed Icebreaker questions",
		},
		model.Command{
			Trigger:          commandIcebreakerShowApproved,
			AutoComplete:     true,
			AutoCompleteDesc: "Show the list of Icebreaker questions",
		},
		model.Command{
			Trigger:          commandIcebreakerResetToDefault,
			AutoComplete:     true,
			AutoCompleteDesc: "Resets the Icebreaker questions to the default ones from this plugin. Channel owners only",
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
		commandIcebreakerApprove: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerApprove(args), nil
		},
		commandIcebreakerReject: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerReject(args), nil
		},
		commandIcebreakerRemove: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerRemove(args), nil
		},
		commandIcebreakerClearAllProposals: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerClearAllProposals(args), nil
		},
		commandIcebreakerClearAllApproved: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerClearAllApproved(args), nil
		},
		commandIcebreakerResetToDefault: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerResetToDefault(args), nil
		},
	}

	userCommands := map[string]func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError){
		commandIcebreakerAdd: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerAdd(args), nil
		},
		commandIcebreakerShowApproved: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerShowApproved(args), nil
		},
		commandIcebreakerShowProposals: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandIcebreakerShowProposals(args), nil
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
	p.FillDefaultQuestions(args.TeamId, args.ChannelId)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("All questions have been reset to the default ones. Beware the pitchforks!"),
	}
}

func (p *Plugin) executeCommandIcebreakerClearAllProposals(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage(args.TeamId, args.ChannelId)
	lenBefore := len(data.ProposedQuestions[args.TeamId][args.ChannelId])
	data.ProposedQuestions[args.TeamId][args.ChannelId] = nil
	p.WriteToStorage(&data)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("All %d proposed questions have been removed. Beware the pitchforks!", lenBefore),
	}
}

func (p *Plugin) executeCommandIcebreakerClearAllApproved(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage(args.TeamId, args.ChannelId)
	lenBefore := len(data.ApprovedQuestions[args.TeamId][args.ChannelId])
	data.ApprovedQuestions[args.TeamId][args.ChannelId] = nil
	p.WriteToStorage(&data)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("All %d questions have been removed. Beware the pitchforks!", lenBefore),
	}
}

func (p *Plugin) executeCommandIcebreakerRemove(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage(args.TeamId, args.ChannelId)

	indeces, errResponse := getIndeces(args.Command, data.ApprovedQuestions[args.TeamId][args.ChannelId])
	if errResponse != nil {
		return errResponse
	}

	for _, index := range indeces {
		//from https://stackoverflow.com/a/37335777/199513
		data.ApprovedQuestions[args.TeamId][args.ChannelId] = append(data.ApprovedQuestions[args.TeamId][args.ChannelId][:index], data.ApprovedQuestions[args.TeamId][args.ChannelId][index+1:]...)
	}
	p.WriteToStorage(&data)

	//TODO: Should we notify the author of the question as well? "Hey <user>! Your question for channel <channel> has been removed: <question>"

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Questions removed",
	}
}

func (p *Plugin) executeCommandIcebreakerReject(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage(args.TeamId, args.ChannelId)

	indeces, errResponse := getIndeces(args.Command, data.ProposedQuestions[args.TeamId][args.ChannelId])
	if errResponse != nil {
		return errResponse
	}

	for _, index := range indeces {
		//from https://stackoverflow.com/a/37335777/199513
		data.ProposedQuestions[args.TeamId][args.ChannelId] = append(data.ProposedQuestions[args.TeamId][args.ChannelId][:index], data.ProposedQuestions[args.TeamId][args.ChannelId][index+1:]...)
	}
	p.WriteToStorage(&data)

	//TODO: Should we notify the author of the proposal as well? "Hey <user>! Your proposed question for channel <channel> has been rejected: <question>"

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Questions rejected",
	}
}

func (p *Plugin) executeCommandIcebreakerApprove(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage(args.TeamId, args.ChannelId)

	indeces, errResponse := getIndeces(args.Command, data.ProposedQuestions[args.TeamId][args.ChannelId])
	if errResponse != nil {
		return errResponse
	}

	for _, index := range indeces {
		question := data.ProposedQuestions[args.TeamId][args.ChannelId][index]
		data.ApprovedQuestions[args.TeamId][args.ChannelId] = append(data.ApprovedQuestions[args.TeamId][args.ChannelId], question)
		//from https://stackoverflow.com/a/37335777/199513
		data.ProposedQuestions[args.TeamId][args.ChannelId] = append(data.ProposedQuestions[args.TeamId][args.ChannelId][:index], data.ProposedQuestions[args.TeamId][args.ChannelId][index+1:]...)
	}

	p.WriteToStorage(&data)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Successfully approved",
	}
}

func (p *Plugin) executeCommandIcebreakerShowProposals(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage(args.TeamId, args.ChannelId)

	if len(data.ProposedQuestions[args.TeamId][args.ChannelId]) == 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "There are no proposed questions for this channel...",
		}
	}

	message := "Proposed questions:\n"
	for index, question := range data.ProposedQuestions[args.TeamId][args.ChannelId] {
		creator := question.Creator
		user, err := p.API.GetUser(creator)
		if err != nil {
			creator = user.GetDisplayName("")
		}
		message = message + fmt.Sprintf("%d.\t%s:\t%s\n", index, creator, question.Question)
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         message,
	}
}

func (p *Plugin) executeCommandIcebreakerShowApproved(args *model.CommandArgs) *model.CommandResponse {
	data := p.ReadFromStorage(args.TeamId, args.ChannelId)

	if len(data.ApprovedQuestions[args.TeamId][args.ChannelId]) == 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "There are no questions for this channel...",
		}
	}

	message := "Questions:\n"
	for index, question := range data.ApprovedQuestions[args.TeamId][args.ChannelId] {
		creator := question.Creator
		user, err := p.API.GetUser(creator)
		if err != nil {
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
	data := p.ReadFromStorage(args.TeamId, args.ChannelId)

	//check if there are any approved questions yet
	if len(data.ApprovedQuestions[args.TeamId][args.ChannelId]) == 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: There are no approved questions that I can ask. Be the first one to propose a question by using '/icebreaker add <question>'",
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
	question := data.ApprovedQuestions[args.TeamId][args.ChannelId][rand.Intn(len(data.ApprovedQuestions[args.TeamId][args.ChannelId]))]
	message := fmt.Sprintf("Hey @%s! %s", user.GetDisplayName(""), question.Question)
	post := &model.Post{
		ChannelId: args.ChannelId,
		RootId:    args.RootId,
		UserId:    p.botID,
		Message:   message,
	}

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

	data := p.ReadFromStorage(args.TeamId, args.ChannelId)

	//Check if the question already is proposed or even approved
	for _, question := range data.ProposedQuestions[args.TeamId][args.ChannelId] {
		if question.Question == newQuestion.Question {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "Error: Your question has already been proposed",
			}
		}
	}
	for _, question := range data.ApprovedQuestions[args.TeamId][args.ChannelId] {
		if question.Question == newQuestion.Question {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "Error: Your question has already been approved",
			}
		}
	}

	data.ProposedQuestions[args.TeamId][args.ChannelId] = append(data.ProposedQuestions[args.TeamId][args.ChannelId], newQuestion)

	p.WriteToStorage(&data)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Thanks %s! Added your proposal: '%s'. Total number of proposals: %d", creator.GetDisplayName(""), newQuestion.Question, len(data.ProposedQuestions[args.TeamId][args.ChannelId])),
	}
}
