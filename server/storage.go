package main

import (
	"bytes"
	"encoding/json"
)

const (
	//KVKEY is the key used for storing the data in the KVStorage
	//KVKEY = "IceBreakerData" //this is the old key which stored the questions in a way more complex way and allowed for proposing and accepting questions by users
	KVKEY = "IceBreakerData_v2"
)

func getDefaultQuestions() []Question {
	//Curated some of the mild questions from https://teambuildinghero.com/icebreaker-questions/
	DefaultQuestions := []Question{
		//mild questions
		Question{Creator: "Icebreaker", Question: "What did you eat for breakfast?"},
		Question{Creator: "Icebreaker", Question: "What is your role in the company?"},
		Question{Creator: "Icebreaker", Question: "What are your favourite pizza toppings?"},
		Question{Creator: "Icebreaker", Question: "What languages do you speak?"},
		Question{Creator: "Icebreaker", Question: "Where were you born?"},
		Question{Creator: "Icebreaker", Question: "Which season is your favourite?"},
		Question{Creator: "Icebreaker", Question: "Do you play any sports?"},
		Question{Creator: "Icebreaker", Question: "Do you have any pets?"},
		Question{Creator: "Icebreaker", Question: "Are you allergic to anything?"},
		Question{Creator: "Icebreaker", Question: "Do you have any Christmas traditions?"},
		Question{Creator: "Icebreaker", Question: "Do you play any musical instruments?"},
		Question{Creator: "Icebreaker", Question: "What was the last movie you attended?"},
		Question{Creator: "Icebreaker", Question: "Where did you grow up?"},
		Question{Creator: "Icebreaker", Question: "Have you ever met a celebrity?"},
		Question{Creator: "Icebreaker", Question: "Do you have any siblings?"},
		Question{Creator: "Icebreaker", Question: "When you were a kid, what did you want to be when you grow up?"},
		Question{Creator: "Icebreaker", Question: "Have you ever broken a bone?"},
		Question{Creator: "Icebreaker", Question: "How many pairs of shoes do you own?"},
		Question{Creator: "Icebreaker", Question: "What is the farthest distance you have driven?"},

		//medium questions
		Question{Creator: "Icebreaker", Question: "What’s your favourite show?"},
		Question{Creator: "Icebreaker", Question: "What book would you recommend other people read?"},
		Question{Creator: "Icebreaker", Question: "Would you prefer luxury beach vacations or backpacking?"},
		Question{Creator: "Icebreaker", Question: "What was your first job?"},
		Question{Creator: "Icebreaker", Question: "What is something you are looking forward to?"},
		Question{Creator: "Icebreaker", Question: "Has your taste in music changed in the last 10 years?"},
		Question{Creator: "Icebreaker", Question: "What is something you do that you don’t like to do?"},
		Question{Creator: "Icebreaker", Question: "What is your favourite fast food restaurant?"},
		Question{Creator: "Icebreaker", Question: "Who is the greatest cook you know?"},
		Question{Creator: "Icebreaker", Question: "Do you have a favourite sports team?"},
		Question{Creator: "Icebreaker", Question: "What game show do you think you could win?"},
		Question{Creator: "Icebreaker", Question: "If you had to move to another country, where would you move?"},
		Question{Creator: "Icebreaker", Question: "What is your favorite Christmas movie?"},
		Question{Creator: "Icebreaker", Question: "What food can you not stand?"},
		Question{Creator: "Icebreaker", Question: "What is the greatest gift you have ever received?"},
		Question{Creator: "Icebreaker", Question: "Are you a dog or cat person?"},
		Question{Creator: "Icebreaker", Question: "Who is your least favorite actor?"},
		Question{Creator: "Icebreaker", Question: "Do you work better in the morning or at night?"},
		Question{Creator: "Icebreaker", Question: "What is something you would like to learn?"},
		Question{Creator: "Icebreaker", Question: "Which is typically better, the book or the movie?"},
		Question{Creator: "Icebreaker", Question: "Do you have a favourite (childhood) video game?"},
		Question{Creator: "Icebreaker", Question: "What is one place you have always wanted to visit?"},
		Question{Creator: "Icebreaker", Question: "If you could drive any car, what car would you drive?"},
		Question{Creator: "Icebreaker", Question: "Do you enjoy rollercoasters?"},
		Question{Creator: "Icebreaker", Question: "What is the best TV show of all time?"},
		Question{Creator: "Icebreaker", Question: "If you could act in any movie, what movie would you be in?"},
		Question{Creator: "Icebreaker", Question: "Say you were given a kitten, what would you name him/her?"},
		Question{Creator: "Icebreaker", Question: "If you could collect anything, what would it be?"},
		Question{Creator: "Icebreaker", Question: "Do you prefer team or individual sports?"},
		Question{Creator: "Icebreaker", Question: "What is the best concert you have ever been to?"},
		Question{Creator: "Icebreaker", Question: "If you could start any business in the world, what would you start?"},
		Question{Creator: "Icebreaker", Question: "Name a piece of technology you wish existed."},
	}

	return DefaultQuestions
}

// FillDefaultQuestions fills in the default questions of this plugin
func (p *Plugin) FillDefaultQuestions() {
	data := p.ReadFromStorage()
	data.Questions = getDefaultQuestions()
	p.WriteToStorage(&data)
}

// ReadFromStorage reads IceBreakerData from the KVStore. Makes sure that data is inited for the given team and channel
func (p *Plugin) ReadFromStorage() IceBreakerData {
	data := IceBreakerData{}
	kvData, err := p.API.KVGet(KVKEY)
	if err != nil {
		//do nothing.. we'll return an empty IceBreakerData then...
	}
	if kvData != nil {
		json.Unmarshal(kvData, &data)
	}

	return data
}

// WriteToStorage writes the given data to storage
func (p *Plugin) WriteToStorage(data *IceBreakerData) {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(data)
	p.API.KVSet(KVKEY, reqBodyBytes.Bytes())
}

// ClearStorage removes all stored data from KVStorage
func (p *Plugin) ClearStorage() {
	p.API.KVDelete(KVKEY)
}
