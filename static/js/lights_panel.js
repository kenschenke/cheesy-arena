// Copyright 2020 Ken Schenke. All Rights Reserved.
// Author: kenschenke@gmail.com (Ken Schenke)
//
// Client-side logic for the field lights panel.

var websocket;

// Sends a websocket message to load a team into an alliance station.
var setFieldLights = function(color)  {
    websocket.send("setFieldLights", color)
};

$(function() {
  // Set up the websocket back to the server.
  websocket = new CheesyWebsocket("/panels/lights/websocket", {});
});
