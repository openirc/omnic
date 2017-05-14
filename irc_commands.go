// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// irc_commands.go - IRC Server commands parser
//

package main

import (
    "log"
    "strings"
)

type Command func(*IRCServer, string)

// Just a dummy function so we don't throw an error
func irc_serv_cmd_untracked(irc_server *IRCServer, message string){}

var commands = map[string]Command {
    "ERROR": irc_serv_cmd_error,
    "NOTICE": irc_serv_cmd_untracked,
    "PING": irc_serv_cmd_ping,
    "001": irc_serv_cmd_rpl_welcome,
    "002": irc_serv_cmd_untracked, // Your host
    "003": irc_serv_cmd_untracked, // Time server was created
    "004": irc_serv_cmd_rpl_myinfo,
    "005": irc_serv_cmd_rpl_capabilities,
    "250": irc_serv_cmd_untracked, // Highest connection count
    "251": irc_serv_cmd_untracked, // Current invisible users + servers
    "252": irc_serv_cmd_untracked, // Current IRC operators
    "254": irc_serv_cmd_untracked, // Current channels
    "255": irc_serv_cmd_untracked, // Current local users
    "266": irc_serv_cmd_untracked, // Current global users
    "372": irc_serv_cmd_untracked, // MOTD
    "375": irc_serv_cmd_untracked, // Start of MOTD
    "376": irc_serv_cmd_untracked, // End of MOTD
    "396": irc_serv_cmd_untracked, // Hidden host
}

func irc_serv_cmd_error(irc_server *IRCServer, message string) {
    log.Printf("Connection error: %s", message)
}

func irc_serv_cmd_ping(irc_server *IRCServer, message string) {
    irc_server.setConnectionState(Disconnected)
    slices := strings.Split(message, " ")
    if len(slices) < 2 {
        log.Printf("Got invalid PING command from server")
        return
    }

    srv := strings.TrimSpace(message[6:])
    cmd := []byte("PONG " + srv + "\n")

    irc_server.sendq <- cmd
}

func irc_serv_cmd_rpl_welcome(irc_server *IRCServer, message string) {
    irc_server.setConnectionState(Connected)
    slices := strings.Split(message, " ")
    if len(slices) < 3 {
        log.Printf("Error: Unable to determine nick")
        return
    }

    nickname := slices[2]
    irc_server.user.setNick(nickname)
}

func irc_serv_cmd_rpl_myinfo(irc_server *IRCServer, message string) {
    slices := strings.Split(message, " ")
    if len(slices) < 5 {
        log.Printf("Error: Unable to determine RPL_MYINFO")
        return
    }

    irc_server.setServerName(slices[3])
    irc_server.setServerVersion(slices[4])
}

func irc_serv_cmd_rpl_capabilities(irc_server *IRCServer, message string) {
    slices := strings.Split(message, " ")
    for idx, capability := range slices {
        // Not important
        if idx < 3 {
            continue
        }

        // We're done
        if capability == ":are" {
            break
        }

        irc_server.setServerCapability(capability)
    }
}

// Parses socket data from the IRCd
func (irc_server *IRCServer) parse_ircd(message string) {
    slices := strings.Split(message, " ")
    if len(slices) == 0 {
        log.Printf("Error: Got 0 slices in ircd message")
        return
    }

    irc_server.logDebugIoIn(message)

    if len(slices) < 3 {
        if method := commands[slices[0]]; method != nil {
            method(irc_server, message)
            return
        }
    }

    cmd := slices[1]

    method := commands[cmd]
    if method != nil {
        method(irc_server, message)
    } else {
        log.Printf("Got unhandled command: %s", cmd)
    }
}

