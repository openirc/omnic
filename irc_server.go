// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// irc_server.go - IRC Server data
//

package main

import (
    "bufio"
    "fmt"
    "net"
    "strconv"
    "strings"
    "time"
)

// IRC Connection state
type ConnectionState int
const (
    Disconnected ConnectionState = iota
    Connecting
    PendingReconnect
    Registering
    Connected
)

type IRCdVerCheck func(*IRCd)

var supported_ircds = map[string]IRCdVerCheck {
    "charybdis": modifyCharybdis,
}

type IRCServer struct {
    // DB settings
    network_id int64
    hostname string
    port int32
    tls bool
    desired_nicks []string
    ident string
    gecos string

    // Websocket clients
    ws_client []*WSClient

    // Connection
    conn net.Conn
    reader *bufio.Reader
    conn_state ConnectionState
    sendq chan[]byte

    // Our client
    user *IRCUser

    // IRCd specifics
    ircd_server_name string   // "prothid.ny.us.gamesurge.net"
    ircd_version string       // "u2.10.12.14(gs1)"
    ircd *IRCd
}

// Setters
func (irc_server *IRCServer) setServerName(name string) {
    irc_server.logDebug(fmt.Sprintf("Setting server name to %s", name))
    irc_server.ircd_server_name = name
}

func (irc_server *IRCServer) setServerVersion(version string) {
    irc_server.logDebug(fmt.Sprintf("Setting server version to %s", version))
    irc_server.ircd_version = version

    for key, f := range supported_ircds {
        if strings.Contains(version, key) {
            f(irc_server.ircd)
            irc_server.logDebug(fmt.Sprintf("Detected IRCd, using %s", irc_server.ircd.ircd_interface))
            return
        }
    }

    irc_server.logDebug(fmt.Sprintf("Unable to detected IRCd, using %s", irc_server.ircd.ircd_interface))
}

func (irc_server *IRCServer) setServerCapability(capability string) {
    if irc_server.ircd == nil {
        panic("Attempted to set server capability when ircd is nil")
    }

    if _, ok := (* irc_server.ircd.capabilities)[capability]; ok {
        irc_server.logDebug(fmt.Sprintf("Enabling capability %s", capability))
        (* irc_server.ircd.capabilities)[capability] = true
    }
}

func (irc_server *IRCServer) setConnectionState(state ConnectionState) {
    irc_server.conn_state = state

    switch state {
    case Connecting:
        irc_server.logDebug("Connecting...")
    case Registering:
        irc_server.logDebug("Registering...")
    case Connected:
        irc_server.logDebug("Connected!")
    case Disconnected:
        irc_server.logDebug("Disconnected.")
    case PendingReconnect:
        irc_server.logDebug("Pending reconnect...")
    default:
        irc_server.logDebug(fmt.Sprintf("Unknown state %d!", state))
    }
}

// Creates a new IRCServer object
func newIRCServer(id int64, hostname string, port int32) *IRCServer {
    return &IRCServer{
        network_id: id,
        hostname: hostname,
        port: port,
        tls: false,
        desired_nicks: make([]string, 0),
        ident: "Omnic",
        gecos: "Omnic",
        ws_client: nil,
        conn: nil,
        sendq: make(chan []byte, 1024),
        user: newIRCUser(),
        ircd: newIRCd(),
    }
}

// Used so we have 1 writer to the IRCServer
// TODO(static): Add throttling
func (irc_server *IRCServer) IRCServerWritePump() {
    for {
        select {
        case message, ok := <-irc_server.sendq:
            // Closed connection
            if !ok {
                irc_server.logDebug("Write pump closed.")
                return
            }

            strmessage := string(message)

            irc_server.logDebugIoOut(strmessage)
            fmt.Fprintf(irc_server.conn, strmessage)

            n := len(irc_server.sendq)
            for i := 0; i < n; i++ {
                msg := <-irc_server.sendq
                strmsg := string(msg)

                irc_server.logDebugIoOut(strmsg)
                fmt.Fprintf(irc_server.conn, strmsg)
            }
        }
    }
}

// Initiailizes the reader and writer. If the reader dies, then we
// close the writer as well.
func (irc_server *IRCServer) connect() {
    for {
        hostname := irc_server.hostname + ":" + strconv.Itoa(int(irc_server.port))
        irc_server.setConnectionState(Connecting)
        conn, err := net.Dial("tcp", hostname)

        if err != nil {
            irc_server.conn = nil
            irc_server.logError(fmt.Sprintf("Unable to connect: %v", err))
            return
        } else {
            irc_server.conn = conn
            irc_server.reader = bufio.NewReader(irc_server.conn)

            irc_server.setConnectionState(Registering)

            go irc_server.IRCServerWritePump()
            irc_server.sendq <- []byte("USER omnic omnic omnic :OpenIRC Project\n")
            irc_server.sendq <- []byte("NICK omnic\n")
        }

        // IRC Socket loop
        for {
            message, buferr := irc_server.reader.ReadString('\n')

            if buferr != nil {
                irc_server.logError(fmt.Sprintf("neterr: %s", buferr))
                break
            } else if len(message) == 0 {
                irc_server.logError(fmt.Sprintf("Lost connection to %s", hostname))
                break
            }
            irc_server.parse_ircd(message)
        }

        irc_server.setConnectionState(PendingReconnect)
        close(irc_server.sendq)
        irc_server.conn = nil

        irc_server.logDebug("Reconnecting in 5 seconds...")

        time.Sleep(5 * time.Second)
    }
}
