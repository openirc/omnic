// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// ws_parser.go - Websocket callback parser
//

package main

import (
    "fmt"
    "strings"
)

func (ws_client *WSClient) parse(bytes_message []byte) {
    message := string(bytes_message)
    parts := strings.Split(message, " ")

    switch parts[0] {
    case "O": // OMNIC
        ws_client.parse_omnic(message, bytes_message)
    case "I": // IRC
        ws_client.parse_irc(message, bytes_message)
    default:
        logError("Got unknown token!")
    }
}

func (ws_client *WSClient) parse_omnic(message string, bytes_message []byte) {
    parts := strings.Split(message, " ")
    switch parts[1] {
    case "A":
        ws_client.try_auth(parts[2])
    default:
        logError("Got unknown omnic command!")
    }
}

func (ws_client *WSClient) parse_irc(message string, bytes_message []byte) {
}

func (ws_client *WSClient) try_auth(token string) {
    logDebug(fmt.Sprintf("Attempting to authenticate %s with token %s\n", ws_client.ip, token))
    auth_id := validate_auth(token)
    logDebug(fmt.Sprintf("Result: %d", auth_id))

    res := fmt.Sprintf("E A %d", auth_id)

    if auth_id > 0 {
        ws_client.setToken(token)
        ws_client.setId(auth_id)
    }

    ws_client.send <- []byte(res)
}
