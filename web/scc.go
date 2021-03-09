// Copyright 2020 Ken Schenke. All Rights Reserved.
// Author: kenschenke@gmail.com (Ken Schenke)
//
// Web routes for SCC interactions.

package web

import (
	"fmt"
	"github.com/Team254/cheesy-arena/field"
	"github.com/Team254/cheesy-arena/model"
	"github.com/Team254/cheesy-arena/websocket"
	"io"
	"log"
	"net/http"
)

// Shows the SCC Testing page.
func (web *Web) sccGetHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}

	template, err := web.parseFiles("templates/setup_scc.html", "templates/base.html")
	if err != nil {
		handleWebErr(w, err)
		return
	}

	data := struct {
		*model.EventSettings
		field.SCCNotifier
	}{
		web.arena.EventSettings,
		web.arena.Scc.GenerateNotifierStatus(),
	}
	err = template.ExecuteTemplate(w, "base", data)
	if err != nil {
		handleWebErr(w, err)
		return
	}
}

// The websocket endpoint for getting realtime updates from the SCC boxes.
func (web *Web) sccWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.NewWebsocket(w, r)
	if err != nil {
		handleWebErr(w, err)
		return
	}
	defer ws.Close()

	// Subscribe the websocket to the notifiers whose messages will be passed on to the client, in a separate goroutine.
	go ws.HandleNotifiers(web.arena.FieldLightsNotifier)
	alliance := ""

	// Loop, waiting for commands and responding to them, until the client closes the connection.
	for {
		messageType, data, err := ws.Read()
		if err != nil {
			if err == io.EOF {
				// Client has closed the connection; nothing to do here.
				if len(alliance) > 0 {
					web.arena.Scc.Disconnect(alliance)
				}
				return
			}
			log.Println(err)
			return
		}

		switch messageType {
		case "sccupdate":
			update := data.(map[string]interface{})
			alliance = ""
			eStop1 := false
			eStop2 := false
			eStop3 := false
			var ok bool
			if alliance, ok = update["alliance"].(string); !ok {
				log.Println("Missing alliance string")
				ws.WriteError("Missing alliance string")
				continue
			}
			if eStop1, ok = update["eStop1"].(bool); !ok {
				log.Println("Missing eStop1 boolean")
				ws.WriteError("Missing eStop1 boolean")
				continue
			}
			if eStop2, ok = update["eStop2"].(bool); !ok {
				log.Println("Missing eStop2 boolean")
				ws.WriteError("Missing eStop2 boolean")
				continue
			}
			if eStop3, ok = update["eStop3"].(bool); !ok {
				log.Println("Missing eStop3 boolean")
				ws.WriteError("Missing eStop3 boolean")
				continue
			}
			web.arena.Scc.ApplyUpdate(field.SCCUpdate{
				Alliance: alliance,
				EStops:[]bool{eStop1,eStop2,eStop3},
			})
		case "setfieldlights":
			color, ok := data.(string)
			if !ok {
				ws.WriteError(fmt.Sprintf("Failed to parse '%s' message.", messageType))
				continue
			}
			switch color {
			case "off":
				web.arena.FieldLights.SetLightsOff(false)
			case "red":
				web.arena.FieldLights.SetLightsRed()
			case "green":
				web.arena.FieldLights.SetLightsGreen()
			case "purple":
				web.arena.FieldLights.SetLightsPurple()
			}
			web.arena.FieldLightsNotifier.Notify()
		default:
			ws.WriteError(fmt.Sprintf("Invalid message type '%s'.", messageType))
			continue
		}
	}
}

// The websocket endpoint for the scc testing page.
func (web *Web) sccGetTestingWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.NewWebsocket(w, r)
	if err != nil {
		handleWebErr(w, err)
		return
	}
	defer ws.Close()

	// Subscribe the websocket to the notifiers whose messages will be passed on to the client, in a separate goroutine.
	go ws.HandleNotifiers(web.arena.SCCNotifier)

	// Loop, waiting for commands and responding to them, until the client closes the connection.
	for {
		messageType, data, err := ws.Read()
		if err != nil {
			if err == io.EOF {
				// Client has closed the connection; nothing to do here.
				return
			}
			log.Println(err)
			return
		}

		switch messageType {
		case "setFieldLights":
			color, ok := data.(string)
			if !ok {
				ws.WriteError(fmt.Sprintf("Failed to parse '%s' message.", messageType))
				continue
			}
			switch color {
			case "off":
				web.arena.FieldLights.SetLightsOff(false)
			case "red":
				web.arena.FieldLights.SetLightsRed()
			case "green":
				web.arena.FieldLights.SetLightsGreen()
			case "purple":
				web.arena.FieldLights.SetLightsPurple()
			}
			web.arena.FieldLightsNotifier.Notify()
		default:
			ws.WriteError(fmt.Sprintf("Invalid message type '%s'.", messageType))
			continue
		}
	}
}

