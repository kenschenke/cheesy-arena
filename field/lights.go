// Copyright 2020 Ken Schenke. All Rights Reserved.
// Author: kenschenke@gmail.com (Ken Schenke)
//
// Functions for controlling field lights

package field

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type LightState int

const (
	LightsOff LightState = iota
	LightsGreen
	LightsRed
	LightsPurple

	lightServer = "http://10.0.100.100:3000/color?color="
)

type Lights struct {
	state LightState
}

type LightApiStatus struct {
	Status string `json:"status"`
	Color string  `json:"color"`
}

func NewLights() (*Lights) {
	lights := new(Lights)
	lights.state = LightsOff

	return lights
}

func (lights *Lights) GetCurrentState() LightState {
	return lights.state
}

func (lights *Lights) GetCurrentStateAsString() string {
	colorStr := "off"
	switch lights.state {
	case LightsOff:
		colorStr = "off"
	case LightsGreen:
		colorStr = "green"
	case LightsRed:
		colorStr = "red"
	case LightsPurple:
		colorStr = "purple"
	}

	return colorStr
}

func (lights *Lights) SetLightsOff() {
	lights.setLights(LightsOff)
}

func (lights *Lights) SetLightsGreen() {
	lights.setLights(LightsGreen)
}

func (lights *Lights) SetLightsRed() {
	lights.setLights(LightsRed)
}

func (lights *Lights) SetLightsPurple() {
	lights.setLights(LightsPurple)
}

func (lights *Lights) setLights(state LightState) {
	if state == lights.state {
		return
	}

	colorStr := "off"
	switch state {
	case LightsOff:
		colorStr = "off"
	case LightsGreen:
		colorStr = "green"
	case LightsRed:
		colorStr = "red"
	case LightsPurple:
		colorStr = "purple"
	}

	url := fmt.Sprintf("%s%s", lightServer, colorStr)
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to set field lights: %s", err.Error())
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to set request to field lights: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed reading response from field lights: %s", err.Error())
		return
	}
	var respData LightApiStatus
	err = json.Unmarshal(body, &respData)
	if err != nil {
		log.Printf("Failed parsing response from field lights: %s", err.Error())
		return
	}

	if respData.Status != "success" {
		log.Printf("Failed to set field lights")
	}

	lights.state = state
}


