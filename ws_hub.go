// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// ws_hub.go - Websocket connection manager
//
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
//

package main

import (
    "log"
)

type WSHub struct {
    clients map[*WSClient]bool
    register chan *WSClient
    unregister chan *WSClient
    irc_servers map[int64]*IRCServer
}

func newWSHub() *WSHub {
    return &WSHub{
        register:     make(chan *WSClient),
        unregister:   make(chan *WSClient),
        clients:      make(map[*WSClient]bool),
        irc_servers:  make(map[int64]*IRCServer),
    }
}

func (h *WSHub) run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            log.Printf("Accepted client: %s", client.ip)
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                log.Printf("Client leaving: %s", client.ip)
                delete(h.clients, client)
                close(client.send)
            }
        }
    }
}
