package main

import (
	"math/rand"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// botID stores the id of our plguin bot
	botID string

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

//Question stores information about a icebreaker question
type Question struct {
	Creator  string `json:"creator"`
	Question string `json:"question"`
}

//IceBreakerData contains all data necessary to be stored for the Icebreaker Plugin
type IceBreakerData struct {
	Questions     []Question `json:"Questions"`
	LastUsers     []string   `json:"LastUsers"`
	LastQuestions []Question `json:"LastQuestions"`
}

//LenHistory sets how many LastUsers/LastQuestions are stored to avoid asking the same users or same questions over and over
const LenHistory int = 50

// OnActivate is invoked when the plugin is activated.
//
// This demo implementation logs a message to the demo channel whenever the plugin is activated.
// It also creates a demo bot account
func (p *Plugin) OnActivate() error {
	//init the rand
	rand.Seed(1337)

	//register all our commands
	if err := p.registerCommands(); err != nil {
		return errors.Wrap(err, "failed to register commands")
	}

	//make sure the bot exists
	botID, ensureBotError := p.Helpers.EnsureBot(&model.Bot{
		Username:    "icebreaker",
		DisplayName: "IceBreaker Bot",
		Description: "A bot created to break the ice",
	}, plugin.ProfileImagePath("/assets/icecube.png"))
	if ensureBotError != nil {
		return errors.Wrap(ensureBotError, "failed to ensure icebreaker bot.")
	}
	p.botID = botID

	return nil
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
