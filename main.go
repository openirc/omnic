// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// main.go - Main file
//

package main

import (
    "log"
    "fmt"
    "flag"
    "net/http"
)

type Profile struct {
    Id int
    Css string
}


func main() {
    // TODO(static): Config
    addr := flag.String("addr", ":8080", "http service address")

    flag.Parse()
    hub := newWSHub()
    go hub.run()
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })

    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }

    fmt.Println("At bottom")
}
