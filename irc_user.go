// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// irc_user.go - IRC User data
//

package main

import (
    "log"
)

type IRCUser struct {
    nick string
    ident string
    host string
    usermodes []byte
    channels []*IRCChannel
    awaymsg string
    isoper bool
}

func newIRCUser() *IRCUser {
    return &IRCUser {
        usermodes: make([]byte, 0),
        channels: make([]*IRCChannel, 0),
    }
}

func (irc_user *IRCUser) setNick(nick string) {
    log.Printf("DEBUG: Setting nickname to %s", nick)
    irc_user.nick = nick
}
