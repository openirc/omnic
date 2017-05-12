// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// main.go - Main file
//

package main

import (
    "log"
    "fmt"
    "flag"
    "net/http"

    "github.com/spf13/viper"
)

type Profile struct {
    Id int
    Css string
}


func main() {
    read_config()
    httpd_bind_ip := viper.GetString("httpd.bindip") + ":" + viper.GetString("httpd.bindport")

    logDebug("Listening on " + httpd_bind_ip)
    addr := flag.String("addr", httpd_bind_ip, "http service address")

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
