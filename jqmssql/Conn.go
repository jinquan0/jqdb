package jqmssql

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/denisenkom/go-mssqldb"
)


func MssqlConn(server string, port int, user string, pass string) (*sql.DB) {
    connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", server, user, pass, port)
    log.Printf(" connString:%s\n", connString)
    db, err := sql.Open("mssql", connString)
    if err != nil {
        log.Fatal("Open connection failed:", err.Error())
    }
    // defer conn.Close()
    return db
}

func MssqlDisconn(MsDb *sql.DB){
    MsDb.Close()
    log.Println("MSSQL/> database disconnected.")
}
