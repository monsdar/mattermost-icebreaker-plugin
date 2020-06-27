package main

import (
	"bytes"
	"encoding/json"
)

const (
	//KVKEY is the key used for storing the data in the KVStorage
	KVKEY = "IceBreakerData"
)

// ReadFromStorage reads IceBreakerData from the KVStore. Makes sure that data is inited for the given team and channel
func (p *Plugin) ReadFromStorage(teamID string, channelID string) IceBreakerData {
	data := IceBreakerData{}
	kvData, err := p.API.KVGet(KVKEY)
	if err != nil {
		//do nothing.. we'll return an empty IceBreakerData then...
	}
	if kvData != nil {
		json.Unmarshal(kvData, &data)
	}

	if data.ProposedQuestions == nil {
		data.ProposedQuestions = make(map[string]map[string][]Question)
	}
	if _, ok := data.ProposedQuestions[teamID]; !ok {
		data.ProposedQuestions[teamID] = make(map[string][]Question)
	}
	if data.ApprovedQuestions == nil {
		data.ApprovedQuestions = make(map[string]map[string][]Question)
	}
	if _, ok := data.ApprovedQuestions[teamID]; !ok {
		data.ApprovedQuestions[teamID] = make(map[string][]Question)
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
