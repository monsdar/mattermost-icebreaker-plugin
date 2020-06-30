package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddIcebreaker(t *testing.T) {
	t.Run("No question given w/o whitespace", func(t *testing.T) {
		plugin := &Plugin{}
		args := &model.CommandArgs{
			Command: "/icebreaker add",
		}
		assert.NotNil(t, plugin.executeCommandIcebreakerAdd(args))
	})

	t.Run("No question given with whitespace", func(t *testing.T) {
		plugin := &Plugin{}
		args := &model.CommandArgs{
			Command: "/icebreaker add ",
		}
		assert.NotNil(t, plugin.executeCommandIcebreakerAdd(args))
	})

	t.Run("Question already proposed", func(t *testing.T) {
		icebreakerData := &IceBreakerData{ProposedQuestions: map[string]map[string][]Question{
			"TestTeam": map[string][]Question{
				"TestChannel": []Question{
					Question{
						Creator: "TestUser", Question: "How do you do?",
					}}}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker add How do you do?",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			UserId:    "TestUser",
		}

		result := plugin.executeCommandIcebreakerAdd(args)
		assert.Equal(t, "Your question has already been proposed", result.Text)
	})

	t.Run("Question already approved", func(t *testing.T) {
		icebreakerData := &IceBreakerData{ApprovedQuestions: map[string]map[string][]Question{
			"TestTeam": map[string][]Question{
				"TestChannel": []Question{
					Question{
						Creator: "TestUser", Question: "How do you do?",
					}}}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker add How do you do?",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			UserId:    "TestUser",
		}

		result := plugin.executeCommandIcebreakerAdd(args)
		assert.Equal(t, "Your question has already been approved", result.Text)
	})

	t.Run("Valid question", func(t *testing.T) {
		icebreakerData := &IceBreakerData{}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		dataAfterAddingTheQuestion := &IceBreakerData{
			ApprovedQuestions: map[string]map[string][]Question{
				"TestTeam": map[string][]Question{}},
			ProposedQuestions: map[string]map[string][]Question{
				"TestTeam": map[string][]Question{
					"TestChannel": []Question{
						Question{
							Creator: "TestUser", Question: "How do you do?",
						}}}}}
		bytesAfterAddingTheQuestion := new(bytes.Buffer)
		json.NewEncoder(bytesAfterAddingTheQuestion).Encode(dataAfterAddingTheQuestion)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", mock.AnythingOfType("string"), bytesAfterAddingTheQuestion.Bytes()).Return(nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker add How do you do?",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			UserId:    "TestUser",
		}

		result := plugin.executeCommandIcebreakerAdd(args)
		assert.Equal(t, "Thanks TestUser! Added your proposal: 'How do you do?'. Total number of proposals: 1", result.Text)
	})
}
