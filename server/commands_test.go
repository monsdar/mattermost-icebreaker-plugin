package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAskIcebreaker_fail(t *testing.T) {
	t.Run("No questions", func(t *testing.T) {
		icebreakerData := &IceBreakerData{Questions: []Question{}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			UserId:    "TestUser",
		}

		result := plugin.executeCommandIcebreaker(args)
		assert.Equal(t, "Error: There are no questions that I can ask. Be the first one to propose a question by using `/icebreaker add <question>`", result.Text)
	})
	t.Run("No users in channel", func(t *testing.T) {
		icebreakerData := &IceBreakerData{Questions: []Question{
			Question{
				Creator: "TestUser", Question: "How do you do?",
			}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return([]*model.User{}, nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			UserId:    "TestUser",
		}

		result := plugin.executeCommandIcebreaker(args)
		assert.Equal(t, "Error: Cannot get a user to ask a question for. Note: This plugin will not ask questions to offline or DND users.", result.Text)
	})
	t.Run("Only bots", func(t *testing.T) {
		icebreakerData := &IceBreakerData{Questions: []Question{
			Question{
				Creator: "TestUser", Question: "How do you do?",
			}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		users := []*model.User{
			&model.User{
				IsBot: true,
			},
			&model.User{
				IsBot: true,
			},
			&model.User{
				IsBot: true,
			},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(users, nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			UserId:    "TestUser",
		}

		result := plugin.executeCommandIcebreaker(args)
		assert.Equal(t, "Error: Cannot get a user to ask a question for. Note: This plugin will not ask questions to offline or DND users.", result.Text)
	})
	t.Run("Only own user", func(t *testing.T) {
		icebreakerData := &IceBreakerData{Questions: []Question{
			Question{
				Creator: "TestUser", Question: "How do you do?",
			}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		users := []*model.User{
			&model.User{
				Id: "TestUser",
			},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(users, nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			UserId:    "TestUser",
		}

		result := plugin.executeCommandIcebreaker(args)
		assert.Equal(t, "Error: Cannot get a user to ask a question for. Note: This plugin will not ask questions to offline or DND users.", result.Text)
	})
	t.Run("Only offline and DND", func(t *testing.T) {
		icebreakerData := &IceBreakerData{Questions: []Question{
			Question{
				Creator: "TestUser", Question: "How do you do?",
			}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		users := []*model.User{
			&model.User{Id: "User1"},
			&model.User{Id: "User2"},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(users, nil)
		api.On("GetUserStatus", "User1").Return(&model.Status{Status: "offline"}, nil)
		api.On("GetUserStatus", "User2").Return(&model.Status{Status: "dnd"}, nil)

		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			UserId:    "TestUser",
		}

		result := plugin.executeCommandIcebreaker(args)
		assert.Equal(t, "Error: Cannot get a user to ask a question for. Note: This plugin will not ask questions to offline or DND users.", result.Text)
	})
}

func TestAskIcebreaker_success(t *testing.T) {
	t.Run("Successful, first user", func(t *testing.T) {
		rand.Seed(1338)
		icebreakerData := &IceBreakerData{Questions: []Question{
			Question{
				Creator: "TestUser", Question: "How do you do?",
			}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		users := []*model.User{
			&model.User{IsBot: true},
			&model.User{Id: "TestUser"},
			&model.User{Id: "User1"},
			&model.User{Id: "User2"},
			&model.User{Id: "SuccessUser", Username: "success_user"},
			&model.User{Id: "SuccessUser2", Username: "success_user2"},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(users, nil)
		api.On("GetUserStatus", "User1").Return(&model.Status{Status: "offline"}, nil)
		api.On("GetUserStatus", "User2").Return(&model.Status{Status: "dnd"}, nil)
		api.On("GetUserStatus", "SuccessUser").Return(&model.Status{Status: "online"}, nil)
		api.On("GetUserStatus", "SuccessUser2").Return(&model.Status{Status: "online"}, nil)
		api.On("CreatePost", &model.Post{
			ChannelId: "TestChannel",
			RootId:    "TestRoot",
			UserId:    "",
			Message:   "Hey @success_user! How do you do?",
		}).Return(nil, nil)

		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker ",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			RootId:    "TestRoot",
			UserId:    "TestUser",
		}

		plugin.ExecuteCommand(nil, args)
	})
	t.Run("Successful, other user", func(t *testing.T) {
		rand.Seed(1337)
		icebreakerData := &IceBreakerData{Questions: []Question{
			Question{
				Creator: "TestUser", Question: "How do you do?",
			}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		users := []*model.User{
			&model.User{IsBot: true},
			&model.User{Id: "TestUser"},
			&model.User{Id: "User1"},
			&model.User{Id: "User2"},
			&model.User{Id: "SuccessUser", Username: "success_user"},
			&model.User{Id: "SuccessUser2", Username: "success_user2"},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(users, nil)
		api.On("GetUserStatus", "User1").Return(&model.Status{Status: "offline"}, nil)
		api.On("GetUserStatus", "User2").Return(&model.Status{Status: "dnd"}, nil)
		api.On("GetUserStatus", "SuccessUser").Return(&model.Status{Status: "online"}, nil)
		api.On("GetUserStatus", "SuccessUser2").Return(&model.Status{Status: "online"}, nil)
		api.On("CreatePost", &model.Post{
			ChannelId: "TestChannel",
			RootId:    "TestRoot",
			UserId:    "",
			Message:   "Hey @success_user2! How do you do?",
		}).Return(nil, nil)

		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			RootId:    "TestRoot",
			UserId:    "TestUser",
		}

		plugin.executeCommandIcebreaker(args)
	})
	t.Run("Successful, history", func(t *testing.T) {
		rand.Seed(1338)
		icebreakerData := &IceBreakerData{Questions: []Question{
			Question{
				Creator: "TestUser", Question: "First question",
			},
			Question{
				Creator: "TestUser", Question: "Second question",
			},
			Question{
				Creator: "TestUser", Question: "Third question",
			}},
			LastUsers: []string{
				"SuccessUser1",
				"SuccessUser2",
			},
			LastQuestions: []Question{
				Question{
					Creator: "TestUser", Question: "First question",
				},
				Question{
					Creator: "TestUser", Question: "Third question",
				},
			},
		}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		users := []*model.User{
			&model.User{IsBot: true},
			&model.User{Id: "TestUser"},
			&model.User{Id: "User1"},
			&model.User{Id: "User2"},
			&model.User{Id: "SuccessUser1", Username: "success_user1"},
			&model.User{Id: "SuccessUser2", Username: "success_user2"},
			&model.User{Id: "SuccessUser3", Username: "success_user3"},
		}

		args := &model.CommandArgs{
			Command:   "/icebreaker ",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			RootId:    "TestRoot",
			UserId:    "TestUser",
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(users, nil)
		api.On("GetUserStatus", "User1").Return(&model.Status{Status: "offline"}, nil)
		api.On("GetUserStatus", "User2").Return(&model.Status{Status: "dnd"}, nil)
		api.On("GetUserStatus", "SuccessUser1").Return(&model.Status{Status: "online"}, nil)
		api.On("GetUserStatus", "SuccessUser2").Return(&model.Status{Status: "online"}, nil)
		api.On("GetUserStatus", "SuccessUser3").Return(&model.Status{Status: "online"}, nil)
		api.On("CreatePost", &model.Post{
			ChannelId: "TestChannel",
			RootId:    "TestRoot",
			UserId:    "",
			Message:   "Hey @success_user3! Second question",
		}).Return(nil, nil)
		plugin.SetAPI(api)

		plugin.ExecuteCommand(nil, args)
	})
}

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

	t.Run("Question too long", func(t *testing.T) {
		plugin := &Plugin{}
		args := &model.CommandArgs{
			Command: "/icebreaker add 123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901",
		}
		assert.NotNil(t, plugin.executeCommandIcebreakerAdd(args))
	})

	t.Run("Question already added", func(t *testing.T) {
		icebreakerData := &IceBreakerData{Questions: []Question{
			Question{
				Creator: "TestUser", Question: "How do you do?",
			}}}
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
		assert.Equal(t, "Error: Your question has already been added", result.Text)
	})

	t.Run("Valid question", func(t *testing.T) {
		icebreakerData := &IceBreakerData{}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		dataAfterAddingTheQuestion := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "How do you do?",
				}}}
		bytesAfterAddingTheQuestion := new(bytes.Buffer)
		json.NewEncoder(bytesAfterAddingTheQuestion).Encode(dataAfterAddingTheQuestion)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser", Id: "TestUserId"}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", "IceBreakerData_v2", bytesAfterAddingTheQuestion.Bytes()).Return(nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command:   "/icebreaker add How do you do?",
			ChannelId: "TestChannel",
			TeamId:    "TestTeam",
			UserId:    "TestUser",
		}

		result := plugin.executeCommandIcebreakerAdd(args)
		assert.Equal(t, "Thanks TestUser! Added your question: 'How do you do?'. Total number of questions: 1", result.Text)
	})
}

func TestRemoveIcebreaker(t *testing.T) {
	t.Run("No index given", func(t *testing.T) {
		icebreakerData := &IceBreakerData{}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		dataAfter := &IceBreakerData{}
		bytesAfter := new(bytes.Buffer)
		json.NewEncoder(bytesAfter).Encode(dataAfter)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser", Id: "TestUserId", Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", "IceBreakerData_v2", bytesAfter.Bytes()).Return(nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command: "/icebreaker admin remove",
		}
		result, _ := plugin.ExecuteCommand(nil, args)
		assert.Equal(t, "Error: Please enter a valid index", result.Text)
	})
	t.Run("Negative index", func(t *testing.T) {
		icebreakerData := &IceBreakerData{}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		dataAfter := &IceBreakerData{}
		bytesAfter := new(bytes.Buffer)
		json.NewEncoder(bytesAfter).Encode(dataAfter)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser", Id: "TestUserId", Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", "IceBreakerData_v2", bytesAfter.Bytes()).Return(nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command: "/icebreaker admin remove -1",
		}
		result, _ := plugin.ExecuteCommand(nil, args)
		assert.Equal(t, "Error: Your given index of -1 is not valid", result.Text)
	})
	t.Run("Index out of bounds", func(t *testing.T) {
		icebreakerData := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "How do you do?",
				}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		dataAfter := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "How do you do?",
				}}}
		bytesAfter := new(bytes.Buffer)
		json.NewEncoder(bytesAfter).Encode(dataAfter)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser", Id: "TestUserId", Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", "IceBreakerData_v2", bytesAfter.Bytes()).Return(nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command: "/icebreaker admin remove 1",
		}
		result, _ := plugin.ExecuteCommand(nil, args)
		assert.Equal(t, "Error: Your given index of 1 is not valid", result.Text)
	})
	t.Run("Success", func(t *testing.T) {
		icebreakerData := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "How do you do?",
				}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		dataAfter := &IceBreakerData{
			Questions: []Question{}}
		bytesAfter := new(bytes.Buffer)
		json.NewEncoder(bytesAfter).Encode(dataAfter)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser", Id: "TestUserId", Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", "IceBreakerData_v2", bytesAfter.Bytes()).Return(nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command: "/icebreaker admin remove 0",
		}
		result, _ := plugin.ExecuteCommand(nil, args)
		assert.Equal(t, "Question removed", result.Text)
	})
	t.Run("Success, multiple questions, remove middle", func(t *testing.T) {
		icebreakerData := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "Index 0",
				},
				Question{
					Creator: "TestUserId", Question: "Index 1",
				},
				Question{
					Creator: "TestUserId", Question: "Index 2",
				},
			}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		dataAfter := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "Index 0",
				},
				Question{
					Creator: "TestUserId", Question: "Index 2",
				},
			}}
		bytesAfter := new(bytes.Buffer)
		json.NewEncoder(bytesAfter).Encode(dataAfter)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser", Id: "TestUserId", Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", "IceBreakerData_v2", bytesAfter.Bytes()).Return(nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command: "/icebreaker admin remove 1",
		}
		result, _ := plugin.ExecuteCommand(nil, args)
		assert.Equal(t, "Question removed", result.Text)
	})
	t.Run("Success, multiple questions, remove end", func(t *testing.T) {
		icebreakerData := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "Index 0",
				},
				Question{
					Creator: "TestUserId", Question: "Index 1",
				},
				Question{
					Creator: "TestUserId", Question: "Index 2",
				},
			}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		dataAfter := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "Index 0",
				},
				Question{
					Creator: "TestUserId", Question: "Index 1",
				},
			}}
		bytesAfter := new(bytes.Buffer)
		json.NewEncoder(bytesAfter).Encode(dataAfter)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser", Id: "TestUserId", Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", "IceBreakerData_v2", bytesAfter.Bytes()).Return(nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command: "/icebreaker admin remove 2",
		}
		result, _ := plugin.ExecuteCommand(nil, args)
		assert.Equal(t, "Question removed", result.Text)
	})
	t.Run("Success, multi index but only first is used", func(t *testing.T) {
		icebreakerData := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "How do you do?",
				}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		dataAfter := &IceBreakerData{
			Questions: []Question{}}
		bytesAfter := new(bytes.Buffer)
		json.NewEncoder(bytesAfter).Encode(dataAfter)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser", Id: "TestUserId", Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", "IceBreakerData_v2", bytesAfter.Bytes()).Return(nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command: "/icebreaker admin remove 0 5",
		}
		result, _ := plugin.ExecuteCommand(nil, args)
		assert.Equal(t, "Question removed", result.Text)
	})
	t.Run("Multi index and first is out of bounds", func(t *testing.T) {
		icebreakerData := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "How do you do?",
				}}}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(icebreakerData)

		dataAfter := &IceBreakerData{
			Questions: []Question{
				Question{
					Creator: "TestUserId", Question: "How do you do?",
				}}}
		bytesAfter := new(bytes.Buffer)
		json.NewEncoder(bytesAfter).Encode(dataAfter)

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "TestUser", Id: "TestUserId", Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("KVSet", "IceBreakerData_v2", bytesAfter.Bytes()).Return(nil)
		plugin.SetAPI(api)

		args := &model.CommandArgs{
			Command: "/icebreaker admin remove 5 0",
		}
		result, _ := plugin.ExecuteCommand(nil, args)
		assert.Equal(t, "Error: Your given index of 5 is not valid", result.Text)
	})
}
