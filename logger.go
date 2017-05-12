// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// logger.go - Logging routines
//

package main

import (
    "log"
    "strings"
)

// XXX: This may be too much abstraction. Feel free to clean
// up.

// IRCServer Logging helpers
func (irc_server *IRCServer) logIRCServerError_(debugType string, message string) {
    log.Printf("ERROR(%s): [%s/%d] %s", debugType, irc_server.hostname, irc_server.network_id, strings.TrimSpace(message))
}

func (irc_server *IRCServer) logIRCServerDebug_(debugType string, message string) {
    log.Printf("DEBUG(%s): [%s/%d] %s", debugType, irc_server.hostname, irc_server.network_id, strings.TrimSpace(message))
}

// IRCServer Logging API.
// TODO(static): Make these variadic functions
func (irc_server *IRCServer) logError(message string) {
    irc_server.logIRCServerError_("Omnic", message)
}

func (irc_server *IRCServer) logDebug(message string) {
    irc_server.logIRCServerDebug_("Omnic", message)
}

func (irc_server *IRCServer) logDebugIoIn(message string) {
    irc_server.logIRCServerDebug_("IO In", message)
}

func (irc_server *IRCServer) logDebugIoOut(message string) {
    irc_server.logIRCServerDebug_("IO Out", message)
}

// Generic Logging API
func logError(message string) {
    log.Printf("ERROR(Omnic): %s", message)
}

func logDebug(message string) {
    log.Printf("DEBUG(Omnic): %s", message)
}
