// Copyright 2020 Ken Schenke. All Rights Reserved.
// Author: kenschenke@gmail.com (Ken Schenke)
//
// Web handlers for field lights interface.

package web

import (
	"fmt"
	"github.com/Team254/cheesy-arena/model"
	"github.com/Team254/cheesy-arena/websocket"
	"io"
	"log"
	"net/http"
)

func (web *Web) lightsPanelHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}

	template, err := web.parseFiles("templates/lights_panel.html", "templates/base.html")
	if err != nil {
		handleWebErr(w, err)
		return
	}

	data := struct {
		*model.EventSettings
	}{web.arena.EventSettings}
	err = template.ExecuteTemplate(w, "base_no_navbar", data)
	if err != nil {
		handleWebErr(w, err)
		return
	}
}

func (web *Web) lightsPanelWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}

	ws, err := websocket.NewWebsocket(w, r)
	if err != nil {
		handleWebErr(w, err)
		return
	}
	defer ws.Close()

	// Loop, waiting for commands and responding to them, until the client closes the connection.
	for {
		command, data, err := ws.Read()
		if err != nil {
			if err == io.EOF {
				// Client has closed the connection; nothing to do here.
				return
			}
			log.Println(err)
			return
		}

		if command == "setFieldLights" {
			color, ok := data.(string)
			if !ok {
				ws.WriteError(fmt.Sprintf("Failed to parse '%s' message.", command))
				continue
			}
			switch color {
			case "off":
				web.arena.FieldLights.SetLightsOff(false)
				web.arena.FieldLightsNotifier.Notify()
			case "red":
				web.arena.FieldLights.SetLightsRed()
				web.arena.FieldLightsNotifier.Notify()
			case "green":
				web.arena.FieldLights.SetLightsGreen()
				web.arena.FieldLightsNotifier.Notify()
			case "purple":
				web.arena.FieldLights.SetLightsPurple()
				web.arena.FieldLightsNotifier.Notify()
			}
		}
	}
}
