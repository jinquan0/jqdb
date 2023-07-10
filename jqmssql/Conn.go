package jqmssql

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/denisenkom/go-mssqldb"
)

// https://github.com/denisenkom/go-mssqldb/blob/master/examples/tsql/tsql.go
func MssqlConn(server string, port int,
                user string, pass string,  database string ) (*sql.DB) {
    dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, pass, port, database)
    log.Printf(" dsn:%s\n", dsn)
    // database driver: Microsoft SQL Server
    db, err := sql.Open("mssql", dsn)
    if err != nil {
        log.Fatal("SQL Server/> Cannot connect: ", err.Error())
        return nil
    }
    defer db.Close()
    err = db.Ping()
    if err != nil {
        log.Fatal("SQL Server/> Cannot connect: ", err.Error())
        return nil
    }
    return db
}

func MssqlDisconn(MsDb *sql.DB){
    MsDb.Close()
    log.Println("SQL Server/> database disconnected.")
}
