// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// ircd_charybdis.go - Charybdis specifics
//

package main

func modifyCharybdis(ircd *IRCd) {
    var usermode_help modeHelpMap
    var chanmode_help modeHelpMap

    for k, v := range default_usermodes_help {
        usermode_help[k] = v
    }

    for k, v := range default_chanmodes_help {
        chanmode_help[k] = v
    }

    chanmode_help['C'] = "Forbid channel CTCPs"
    chanmode_help['S'] = "SSL Only"

    ircd.usermode_help = &usermode_help
    ircd.chanmode_help = &chanmode_help
    ircd.ircd_interface = "Charybdis"
}
