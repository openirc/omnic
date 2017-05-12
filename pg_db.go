// Copyright 2017 OpenIRC - https://github.com/openirc/omnic
// pg_db.go - PostgreSQL database routines
//

package main

import (
    "fmt"

    "database/sql"
    _ "github.com/lib/pq"
)

// TODO(static): Configize this
func get_sql() *sql.DB {
    db, err := sql.Open("postgres", "dbname=openirc")

    if err != nil {
        fmt.Println("Postgres error:", err)
        return nil
    }

    return db
}

func validate_auth(token string) int64 {
    db := get_sql()
    if db == nil {
        return -1
    }

    rows, err := db.Query("SELECT user_id FROM auth_tokens WHERE token $1", token)

    if err != nil {
        fmt.Println("Failed to query:", err)
        return -1
    }

    var ids []int64

    for rows.Next() {
        var user_id int64
        if scan_err := rows.Scan(&user_id); scan_err != nil {
            fmt.Println("Failed to scan", scan_err)
        } else {
            ids = append(ids, user_id)
        }
    }

    if len(ids) == 0 {
        return 0
    } else if len(ids) > 1 {
        fmt.Println("Got too many tokens!")
        return 0
    }

    return ids[0]

}
