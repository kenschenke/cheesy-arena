// Copyright 2020 Ken Schenke. All Rights Reserved.
// Author: kenschenke@gmail.com (Ken Schenke)
//
// Functions for controlling field lights

package field

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

	lights.state = state
}


