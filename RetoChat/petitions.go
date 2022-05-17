package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type printChannel struct {
	ChannelName string `json:"channelName"`
	ClientName  string `json:"clientName"`
}

type printStatistic struct {
	Name      string `json:"name"`
	Y         int    `json:"y"`
	Drilldown string `json:"drilldown"`
}

func (s *server) Start(w http.ResponseWriter, r *http.Request) {

	channel := printChannel{}
	arrayChannels := []printChannel{}

	for k := range s.channels {
		for k2 := range s.channels[k].members {
			var nameChannel, nameClient string

			nameChannel = k
			nameClient = s.channels[k].members[k2].name

			channel.ChannelName = nameChannel
			channel.ClientName = nameClient

			arrayChannels = append(arrayChannels, channel)

		}
	}

	fileJson, _ := json.MarshalIndent(arrayChannels, "", " ")
	_ = ioutil.WriteFile("web/channels.json", fileJson, 0644)

}

func (s *server) Statistic(w http.ResponseWriter, r *http.Request) {

	statistic := printStatistic{}
	arrayStatistic := []printStatistic{}

	for k := range s.channels {
		var count1 int = 0
		for range s.channels[k].sending {
			count1++
		}
		statistic.Name = k
		statistic.Y = count1
		statistic.Drilldown = k

		arrayStatistic = append(arrayStatistic, statistic)
	}

	fileJson, _ := json.Marshal(arrayStatistic)
	_ = ioutil.WriteFile("web/statistic.json", fileJson, 0644)

}
