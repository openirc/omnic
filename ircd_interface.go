// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// ircd_interface.go - Abstract IRCd interface
//

package main

// IRCd capabilities
type CapabilitiesMap map[string]bool
var capabilities = CapabilitiesMap {
    "CNOTICE": false,    // Channel notice
    "CPRIVMSG": false,   // Channel privmsg
    "ETRACE": false,     // Don't remember
    "EXCEPTS": false,    // Ban Exceptions +e (I think)
    "FNC": false,        // force nick change (SANICK/SVSNICK/etc)
    "INVEX": false,      // Invite Exceptions +I
    "KNOCK": false,      // Don't remember
    "SAFELIST": false,   // Don't remember
    "WHOX": false,       // WHO #chan %flags
}

// Mode help menus (Maybe this should all be done client side...)
type modeHelpMap map[byte]string
var default_usermodes_help = modeHelpMap {
    'o': "IRC Operator",
}

var default_chanmodes_help = modeHelpMap {
    's': "Secret channel",
}

// Basic IRCd
type IRCd struct {
    capabilities *CapabilitiesMap   // "WHOX INVEX FNC"
    channel_types []byte            // "#&"
    cmode_prefix []byte             // "~&@%v"
    chan_opmodes []byte             // "qaohv"
    max_topiclen int                // 256
    max_modechange int              // 4
    list_modes []byte               // "be"
    param_modes []byte              // "k"
    param_modes_when_set []byte     // "l"
    plain_modes []byte              // "pstnmi"
    supports_whox bool              // If we support WHOX, cache this
    usermode_help *modeHelpMap  // Help list of user modes
    chanmode_help *modeHelpMap  // Help list of channel modes
    ircd_interface string           // What data we're pulling from
}

func newIRCd() *IRCd {
    return &IRCd{
        capabilities: &capabilities,
        channel_types: make([]byte, 0),
        cmode_prefix: make([]byte, 0),
        chan_opmodes: make([]byte, 0),
        max_topiclen: 256,
        list_modes: make([]byte, 0),
        param_modes: make([]byte, 0),
        param_modes_when_set: make([]byte, 0),
        plain_modes: make([]byte, 0),
        usermode_help: &default_usermodes_help,
        chanmode_help: &default_chanmodes_help,
        ircd_interface: "Generic IRCd",
    }
}
